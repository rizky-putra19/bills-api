package internal

import (
	"gitlab.com/lokalpay-dev/digital-goods/internal/dto"
	"gitlab.com/lokalpay-dev/digital-goods/internal/entity"
)

type ProductRepositoryItf interface {
	GetProduct(productID int) (*entity.Product, error)
	GetProducts() ([]entity.Product, error)
}

type PaymentRepositoryItf interface {
	CreatePayment(p *entity.Payment) error
	GetPaymentByReferenceNumber(referenceNumber string) (*entity.Payment, error)
	GetPayment(id int) (*entity.Payment, error)
	UpdatePayment(p *entity.Payment) error
}

type OrderRepositoryItf interface {
	CreateOrder(o *entity.Order) (int, error)
	GetOrder(id int) (*entity.Order, error)
	GetOrderDetails(id int) (dto.OrderDetails, error)
	GetOrderByCustomer(id int) ([]dto.OrderDetails, error)
	UpdateOrder(o *entity.Order) error
}

type OrderItemsRepositoryItf interface {
	CreateOrderItem(oi *entity.OrderItem) error
	GetOrderItem(id int) (*entity.OrderItem, error)
	GetOrderItemByOrderID(orderID int) (*entity.OrderItem, error)
	GetOrderItemByReferenceNumber(referenceNumber string) (*entity.OrderItem, error)
	UpdateOrderItem(oi *entity.OrderItem) error
}

type CustomerRepositoryItf interface {
	CreateCustomer(email, passwordHash string) (int, error)
	GetCustomerByEmail(email string) (entity.Customer, error)
	GetCustomerByID(id int) (entity.Customer, error)
}
