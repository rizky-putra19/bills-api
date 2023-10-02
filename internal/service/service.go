package service

import (
	"gitlab.com/lokalpay-dev/digital-goods/config"
	"gitlab.com/lokalpay-dev/digital-goods/internal"
	"gitlab.com/lokalpay-dev/digital-goods/internal/repository"
)

type Service struct {
	Product           internal.ProductServiceItf
	Order             internal.OrderServiceItf
	Payment           internal.PaymentServiceItf
	Fulfillment       internal.FulfillmentServiceItf
	Customer          internal.CustomerServiceItf
	AdapterLFI        internal.LFIProvider
	AdapterPrismalink internal.PrismalinkProvider
}

// New is dependency injection for service layer
func New(
	repo *repository.Repository,
	cfg config.App,
	adapterLFI internal.LFIProvider,
	adapterPrismalink internal.PrismalinkProvider,
) *Service {
	product := NewProduct(repo.Product)
	order := NewOrder(
		repo.Product,
		repo.Order,
		repo.Payment,
		repo.OrderItem,
		adapterLFI,
	)
	payment := NewPayment(repo.Payment)
	fulfillment := NewFulfillment(
		repo.Order,
		repo.OrderItem,
		repo.Payment,
		product,
		repo.Product,
		adapterLFI,
		adapterPrismalink,
	)
	customer := NewCustomer(repo.Customer, cfg)

	return &Service{
		Product:     product,
		Order:       order,
		Payment:     payment,
		Fulfillment: fulfillment,
		Customer:    customer,
	}
}
