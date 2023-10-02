package service

import (
	"errors"

	"gitlab.com/lokalpay-dev/digital-goods/internal"
	"gitlab.com/lokalpay-dev/digital-goods/internal/constant"
	"gitlab.com/lokalpay-dev/digital-goods/internal/dto"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/helper"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/slog"
)

type Fulfillment struct {
	orderRepo         internal.OrderRepositoryItf
	orderItemRepo     internal.OrderItemsRepositoryItf
	paymentRepo       internal.PaymentRepositoryItf
	productService    internal.ProductServiceItf
	productRepo       internal.ProductRepositoryItf
	adapterLfi        internal.LFIProvider
	adapterPrismalink internal.PrismalinkProvider
}

func NewFulfillment(
	orderRepo internal.OrderRepositoryItf,
	orderItemRepo internal.OrderItemsRepositoryItf,
	paymentRepo internal.PaymentRepositoryItf,
	productService internal.ProductServiceItf,
	productRepo internal.ProductRepositoryItf,
	adapterLfi internal.LFIProvider,
	adapterPrismalink internal.PrismalinkProvider,
) *Fulfillment {
	return &Fulfillment{
		orderRepo:         orderRepo,
		orderItemRepo:     orderItemRepo,
		paymentRepo:       paymentRepo,
		productService:    productService,
		productRepo:       productRepo,
		adapterLfi:        adapterLfi,
		adapterPrismalink: adapterPrismalink,
	}
}

func (f *Fulfillment) OrderFulfillment(orderID int) error {
	/**
	- lookup into order_items
	- send request to tokopedia
	- if success then set order_items.fulfillment status to processed
	*/
	order, err := f.orderRepo.GetOrder(orderID)
	if err != nil {
		slog.Errorw(
			"fulfillment failed to get order details",
			"order_id", orderID,
		)
		return err
	}

	if order.Status != constant.OrderStatusCreated {
		slog.Errorw(
			"fulfillment failed order has been processed",
			"order_id", orderID,
		)
		return errors.New("fulfillment failed order has been processed")

	}

	order.Status = constant.OrderStatusPending
	err = f.orderRepo.UpdateOrder(order)
	if err != nil {
		slog.Errorw(
			"fulfillment failed to set order status into pending",
			"order_id", orderID,
		)
		return err
	}

	_, err = f.orderItemRepo.GetOrderItemByOrderID(orderID)
	if err != nil {
		slog.Errorw(
			"fulfillment failed to get order item",
			"order_id", orderID,
		)
		return err
	}

	// TODO: send request to Tokopedia endpoint to order selected item and get reference number

	// TODO: set order_item.fulfillment_status to pending and set order_item.fulfillment_reference_number to tokopedia_reference_number

	return nil
}

func (f *Fulfillment) OrderFulfillmentPrismalinkProvider(referenceNumber string) (dto.FullfilmentResponse, error) {
	/**
	- check status on LFI
	- if success, lookup to payments set status success
	- send request to biller Prismalink
	- if success, set success order.status and order_items.fulfillment_status
	*/

	var fulfillmentResponse dto.FullfilmentResponse

	// check status on LFI
	paymentResponse, err := f.adapterLfi.SyncStatus(referenceNumber)
	if err != nil {
		slog.Errorw(
			"fulfillment failed to sync status with LFI provider",
			"reference_number", referenceNumber,
		)
		return fulfillmentResponse, err
	}

	// lookup to payments and update
	paymentEntity, err := f.paymentRepo.GetPaymentByReferenceNumber(referenceNumber)
	if err != nil {
		slog.Errorw(
			"fulfillment failed to get payment",
			"reference_number", referenceNumber,
		)
		return fulfillmentResponse, err
	}

	if paymentResponse.Process.Status == constant.PaymentStatusExpired {
		slog.Errorw(
			"fulfillment failed get status not success from LFI provider",
			"reference_number", referenceNumber,
		)

		// update payment status failed
		paymentEntity.Status = constant.PaymentStatusExpired
		paymentEntity.AccountNumber = paymentResponse.PaymentMethodResponse.Account.AccountNumber
		_ = f.paymentRepo.UpdatePayment(paymentEntity)

		return fulfillmentResponse, errors.New("payment failed")
	}

	if paymentResponse.Process.Status == constant.PaymentStatusCreated {
		slog.Errorw(
			"fulfillment failed get status pending from LFI provider",
			"reference_number", referenceNumber,
		)

		return fulfillmentResponse, errors.New("payment pending")
	}

	if paymentEntity.Status != constant.PaymentStatusCreated {
		slog.Errorw(
			"fulfillment has been proceed",
			"reference_number", referenceNumber,
		)

		return fulfillmentResponse, errors.New("payment has been proceed")
	}

	// update payment status success
	paymentEntity.Status = constant.PaymentStatusSuccess
	paymentEntity.AccountNumber = paymentResponse.PaymentMethodResponse.Account.AccountNumber
	err = f.paymentRepo.UpdatePayment(paymentEntity)
	if err != nil {
		slog.Errorw(
			"fulfilment failed to update payment status",
			"reference_number", referenceNumber,
		)
		return fulfillmentResponse, err
	}

	// lookup order_items and orders
	orderItemEntity, err := f.orderItemRepo.GetOrderItemByReferenceNumber(referenceNumber)
	if err != nil {
		slog.Errorw(
			"fulfillment failed to get order item entity",
			"reference_number", referenceNumber,
		)
		return fulfillmentResponse, err
	}

	orderEntity, err := f.orderRepo.GetOrder(orderItemEntity.OrderID)
	if err != nil {
		slog.Errorw(
			"fulfillment failed to get order entity",
			"reference_number", referenceNumber,
		)
		return fulfillmentResponse, err
	}

	// request payment to biller prismalink
	product, err := f.productRepo.GetProduct(orderItemEntity.ProductID)
	if err != nil {
		slog.Errorw(
			"fulfillment failed to get product entity",
			"reference_number", referenceNumber,
		)
		return fulfillmentResponse, err
	}

	prefixPhoneNumber := helper.GetInputNumber(orderItemEntity.ClientNumber)
	instCd := f.productService.OperatorSelection(prefixPhoneNumber)
	billerRequestPayload := dto.RequestPayloadPrismalink{
		IdTransaction: referenceNumber,
		BillNo:        orderItemEntity.ClientNumber,
		InstCd:        instCd,
		Amount:        constant.PrismalinkPrePaidProduct[product.ProductName],
	}

	// request to prismalink
	reqType := constant.PrismalinkPaymentRequest
	billerResponse, err := f.adapterPrismalink.RequestBiller(billerRequestPayload, reqType)
	if err != nil {
		slog.Errorw(
			"fulfillment failed to request biller",
			"referece_number", referenceNumber,
		)
		return fulfillmentResponse, err
	}

	if billerResponse.Rc != constant.PrismalinkBillerResponseSuccess {
		slog.Errorw(
			"fulfillment failed biller response not success",
			"referece_number", referenceNumber,
		)

		// update orders and order_items status
		orderItemEntity.FulfillmentStatus = constant.FulfillmentStatusFailed
		orderEntity.Status = constant.OrderStatusFailed

		_ = f.orderItemRepo.UpdateOrderItem(orderItemEntity)
		_ = f.orderRepo.UpdateOrder(orderEntity)

		return fulfillmentResponse, err
	}

	// update success order_items entity
	orderItemEntity.FulfillmentStatus = constant.FulfillmentStatusSuccess
	orderEntity.Status = constant.OrderStatusSuccess

	_ = f.orderItemRepo.UpdateOrderItem(orderItemEntity)
	_ = f.orderRepo.UpdateOrder(orderEntity)

	fulfillmentResponse = dto.FullfilmentResponse{
		OrderId:     orderEntity.ID,
		OrderStatus: orderEntity.Status,
	}

	return fulfillmentResponse, nil
}

func (f *Fulfillment) SyncFulfillmentStatus(fulfillmentReferenceNumber string) error {

	/**
	- capture request from tokopedia < this function is on controller layer
	- lookup into order_items
	- if status is pending then set order_items.fulfillment status to success
	- maybe notify user
	*/
	orderItem, err := f.orderItemRepo.GetOrderItemByReferenceNumber(fulfillmentReferenceNumber)
	if err != nil {
		slog.Errorw(
			"fulfillment failed to get order item",
			"fulfillment_reference_number", fulfillmentReferenceNumber,
		)
		return err
	}

	if orderItem.FulfillmentStatus != constant.FulfillmentStatusPending {
		slog.Errorw(
			"order items is not on pending state",
			"fulfillment_reference_number", fulfillmentReferenceNumber,
		)
		return errors.New("order items is not on pending state")
	}

	orderItem.FulfillmentStatus = constant.FulfillmentStatusSuccess
	err = f.orderItemRepo.UpdateOrderItem(orderItem)
	if err != nil {
		slog.Errorw(
			"failed to set fulfillment status into success",
			"fulfillment_reference_number", fulfillmentReferenceNumber,
		)
		return err
	}

	return nil
}
