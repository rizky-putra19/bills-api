package psql

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"gitlab.com/lokalpay-dev/digital-goods/internal/dto"
	"gitlab.com/lokalpay-dev/digital-goods/internal/entity"
)

type Order struct {
	db *sqlx.DB
}

func NewOrder(db *sqlx.DB) *Order {
	return &Order{
		db: db,
	}
}

// CreateOrder inserts a new order into the database
func (or *Order) CreateOrder(o *entity.Order) (int, error) {
	query := "INSERT INTO orders (customer_id, order_date, total_price, status, created_at) VALUES ($1, $2, $3, $4, now()) RETURNING id"
	row := or.db.QueryRow(query, o.CustomerID, o.OrderDate, o.TotalPrice, o.Status)
	err := row.Scan(&o.ID)
	if err != nil {
		return 0, err
	}
	return o.ID, nil
}

// GetOrder retrieves an order from the database by ID
func (or *Order) GetOrder(id int) (*entity.Order, error) {
	o := &entity.Order{}
	query := "SELECT * FROM orders WHERE id=$1"
	err := or.db.Get(o, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return o, nil
}

func (or *Order) UpdateOrder(o *entity.Order) error {
	query := "UPDATE orders SET customer_id=$1, order_date=$2, total_price=$3, status=$4, updated_at=now() WHERE id=$5"
	_, err := or.db.Exec(query, o.CustomerID, o.OrderDate, o.TotalPrice, o.Status, o.ID)
	if err != nil {
		return err
	}
	return nil
}

func (or *Order) GetOrderDetails(id int) (dto.OrderDetails, error) {
	var o dto.OrderDetails
	query := `
		select
			o.id,
			o.status as "order_status",
			o.customer_id,
			p.status as "payment_status",
			o.created_at as "order_created_at",
			o.total_price,
			oi.client_number,
			pr.product_name,
			oi.fulfillment_status,
			COALESCE(p.account_number,'') as "account_number",
			p.payment_method,
			p.payment_expiry_time
		from orders o
		left join payments p on o.id = p.order_id
		left join order_items oi on oi.order_id = o.id
		left join products pr on pr.id = oi.product_id
		where o.id = $1
	`

	err := or.db.Get(&o, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.OrderDetails{}, errors.New("order not found")
		}
		return dto.OrderDetails{}, err
	}
	return o, nil
}

func (or *Order) GetOrderByCustomer(id int) ([]dto.OrderDetails, error) {
	var o []dto.OrderDetails
	query := `
		select
			o.id,
			o.status as "order_status",
			o.customer_id,
			p.status as "payment_status",
			o.created_at as "order_created_at",
			o.total_price,
			oi.client_number,
			pr.product_name,
			oi.fulfillment_status,
			COALESCE(p.account_number,'') as "account_number",
			p.payment_method,
			p.payment_expiry_time
		from orders o
		left join payments p on o.id = p.order_id
		left join order_items oi on oi.order_id = o.id
		left join products pr on pr.id = oi.product_id
		where o.customer_id = $1
	`

	err := or.db.Select(&o, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []dto.OrderDetails{}, errors.New("order not found")
		}
		return []dto.OrderDetails{}, err
	}
	return o, nil
}
