package service

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"gitlab.com/lokalpay-dev/digital-goods/internal"
	"gitlab.com/lokalpay-dev/digital-goods/internal/constant"
	"gitlab.com/lokalpay-dev/digital-goods/internal/dto"
	"gitlab.com/lokalpay-dev/digital-goods/internal/entity"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/slog"
)

type Order struct {
	regex         *regexp.Regexp
	productRepo   internal.ProductRepositoryItf
	orderRepo     internal.OrderRepositoryItf
	paymentRepo   internal.PaymentRepositoryItf
	orderItemRepo internal.OrderItemsRepositoryItf
	adapterLFI    internal.LFIProvider
}

func NewOrder(
	productRepo internal.ProductRepositoryItf,
	orderRepo internal.OrderRepositoryItf,
	paymentRepo internal.PaymentRepositoryItf,
	orderItemRepo internal.OrderItemsRepositoryItf,
	adapterLFI internal.LFIProvider,
) *Order {
	// regex only allow string
	reg, _ := regexp.Compile("[^a-zA-Z]+")
	return &Order{
		regex:         reg,
		productRepo:   productRepo,
		orderRepo:     orderRepo,
		paymentRepo:   paymentRepo,
		orderItemRepo: orderItemRepo,
		adapterLFI:    adapterLFI,
	}
}

func (o *Order) FindOrderDetails(id int) (dto.OrderDetails, error) {
	order, err := o.orderRepo.GetOrderDetails(id)
	if err != nil {
		return dto.OrderDetails{}, err
	}
	return order, nil
}

func (o *Order) FindOrderByCustomer(customerID int) ([]dto.OrderDetails, error) {
	order, err := o.orderRepo.GetOrderByCustomer(customerID)
	if err != nil {
		return []dto.OrderDetails{}, err
	}
	return order, nil
}

func (o *Order) SaveOrder(payload dto.OrderPayload) (dto.PaymentMetadata, error) {
	/**
	- select product get product price
	- insert into orders
	- insert into order items
	- insert into payment
	*/
	now := time.Now()
	product, err := o.productRepo.GetProduct(payload.ProductID)
	if err != nil {
		slog.Errorw("order request failed when getting product details", "err", err.Error())
		return dto.PaymentMetadata{}, err
	}

	totalPrice := product.COGS + product.PartnerFee

	orderData := &entity.Order{
		CustomerID: payload.CustomerID,
		OrderDate:  now,
		TotalPrice: totalPrice,
		Status:     constant.OrderStatusCreated,
		CreatedAt:  &now,
	}

	orderID, err := o.orderRepo.CreateOrder(orderData)
	if err != nil {
		slog.Errorw("order request failed when creating new order", "err", err.Error())
		return dto.PaymentMetadata{}, err
	}

	// reference number for fulfillment reference number
	refNo := o.generateReferenceNumber()

	orderItemData := &entity.OrderItem{
		OrderID:                    orderID,
		ProductID:                  payload.ProductID,
		ClientNumber:               payload.ClientNumber,
		Price:                      totalPrice,
		COGS:                       product.COGS,
		FulfillmentStatus:          constant.FulfillmentStatusCreated,
		FulfillmentReferenceNumber: refNo,
		CreatedAt:                  &now,
	}
	err = o.orderItemRepo.CreateOrderItem(orderItemData)
	if err != nil {
		slog.Errorw("order request failed when insert order item", "err", err.Error())
		return dto.PaymentMetadata{}, err
	}

	oneHour := now.Add(time.Hour)

	// request token for payment
	token, err := o.adapterLFI.RequestAuthToken(refNo)
	if err != nil {
		slog.Errorw("order request failed when request token to adapter LFI", "err", err.Error())
		return dto.PaymentMetadata{}, err
	}

	paymentData := &entity.Payment{
		OrderID:           orderID,
		PaymentDate:       nil,
		PaymentMethod:     "virtual-account",
		PaymentExpiryTime: &oneHour,
		PaymentAmount:     totalPrice + 1000, // provider admin fee
		AdminFee:          1000,
		Status:            constant.PaymentStatusCreated,
		ReferenceNumber:   refNo, // generate reference_number
		AccountNumber:     "",
	}
	err = o.paymentRepo.CreatePayment(paymentData)
	if err != nil {
		slog.Errorw("order request failed when payment creation", "err", err.Error())
		return dto.PaymentMetadata{}, err
	}

	return dto.PaymentMetadata{
		OrderID:           orderID,
		ReferenceNumber:   paymentData.ReferenceNumber,
		AccountNumber:     paymentData.AccountNumber,
		Amount:            paymentData.PaymentAmount,
		LFIToken:          token,
		PaymentExpiryTime: oneHour.Format(time.RFC3339),
	}, err
}

func (o *Order) generateVABasedOnBank(bankCode string) string {
	// please put min 15 character because lower than that will cause intermittent issue
	min := 100000000000000
	max := 999999999999999

	randNumber := rand.Intn((max - min + 1) + min)

	prefixBank := constant.TransformToVAPrefix[bankCode]
	generatedNumber := fmt.Sprintf("%v%d", prefixBank, randNumber)
	return generatedNumber[0:16]
}

func (o *Order) generateReferenceNumber() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const segmentLength = 7 // the length of each first and second segment
	rand.Seed(time.Now().UnixNano())

	// generate 2 segments with a length of 7 characters
	segment1 := make([]byte, segmentLength)
	for i := range segment1 {
		segment1[i] = charset[rand.Intn(len(charset))]
	}
	segment2 := make([]byte, segmentLength)
	for i := range segment2 {
		segment2[i] = charset[rand.Intn(len(charset))]
	}

	return "LPID-" + string(segment1) + "-" + string(segment2)
}
