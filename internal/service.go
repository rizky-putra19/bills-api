package internal

import (
	"gitlab.com/lokalpay-dev/digital-goods/internal/dto"
	"gitlab.com/lokalpay-dev/digital-goods/internal/entity"
)

//go:generate mockgen -destination=service_mock.go -package=internal -source=service.go

type ProductServiceItf interface {
	GetProducts() ([]entity.Product, error)
	OperatorSelection(inputNumber string) string
}

type PaymentServiceItf interface {
	// AcceptPayment ...
	AcceptPayment(referenceNumber string, amount float64) (int, error)
}

type OrderServiceItf interface {
	SaveOrder(payload dto.OrderPayload) (dto.PaymentMetadata, error)
	FindOrderDetails(id int) (dto.OrderDetails, error)
	FindOrderByCustomer(customerID int) ([]dto.OrderDetails, error)
}

type FulfillmentServiceItf interface {
	OrderFulfillment(orderId int) error
	OrderFulfillmentPrismalinkProvider(referenceNumber string) (dto.FullfilmentResponse, error)
	SyncFulfillmentStatus(fulfillmentReferenceNumber string) error
}

type CustomerServiceItf interface {
	Registration(payload dto.AuthPayload) (dto.AuthResponse, error)
	Authentication(payload dto.AuthPayload) (dto.AuthResponse, error)
	GetProfile(id int) (entity.Customer, error)
}
