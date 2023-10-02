package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"gitlab.com/lokalpay-dev/digital-goods/config"
	"gitlab.com/lokalpay-dev/digital-goods/internal"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/slog"
	"gitlab.com/lokalpay-dev/digital-goods/internal/repository/psql"
)

type Repository struct {
	db        *sqlx.DB
	Payment   internal.PaymentRepositoryItf
	Product   internal.ProductRepositoryItf
	Order     internal.OrderRepositoryItf
	OrderItem internal.OrderItemsRepositoryItf
	Customer  internal.CustomerRepositoryItf
}

func New(cfg config.Storage) *Repository {
	dbDriver, _ := initializeConnection(cfg.PSQL)
	slog.Infow("sql connection open", "dbname", cfg.PSQL.DBName)

	customer := psql.NewCustomer(dbDriver)
	payment := psql.NewPayment(dbDriver)
	product := psql.NewProduct(dbDriver)
	order := psql.NewOrder(dbDriver)
	orderItem := psql.NewOrderItems(dbDriver)

	return &Repository{
		db:        dbDriver,
		Customer:  customer,
		Product:   product,
		Payment:   payment,
		Order:     order,
		OrderItem: orderItem,
	}
}

func initializeConnection(psql config.PSQL) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		psql.Host,
		psql.Port,
		psql.User,
		psql.Password,
		psql.DBName,
	)
	db, err := sqlx.Open("postgres", connectionString)
	// defer db.Close()
	if err != nil {
		slog.Fatalw(
			"failed to initialize connection",
			"error",
			err.Error(),
		)
	}

	return db, err
}
