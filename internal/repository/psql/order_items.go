package psql

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"gitlab.com/lokalpay-dev/digital-goods/internal/entity"
)

type OrderItems struct {
	db *sqlx.DB
}

func NewOrderItems(db *sqlx.DB) *OrderItems {
	return &OrderItems{
		db: db,
	}
}

// CreateOrderItem inserts a new order item into the database
func (oit *OrderItems) CreateOrderItem(oi *entity.OrderItem) error {
	query := "INSERT INTO order_items (order_id, product_id, price, cogs, client_number, fulfillment_status, fulfillment_reference_number) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	row := oit.db.QueryRow(query, oi.OrderID, oi.ProductID, oi.Price, oi.COGS, oi.ClientNumber, oi.FulfillmentStatus, oi.FulfillmentReferenceNumber)
	err := row.Scan(&oi.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetOrderItem retrieves an order item from the database by ID
func (oit *OrderItems) GetOrderItem(id int) (*entity.OrderItem, error) {
	oi := &entity.OrderItem{}
	query := "SELECT * FROM order_items WHERE id=$1"
	err := oit.db.Get(oi, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order item not found")
		}
		return nil, err
	}
	return oi, nil
}

// GetOrderItemByOrderID retrieves an order item from the database by order_id
func (oit *OrderItems) GetOrderItemByOrderID(orderID int) (*entity.OrderItem, error) {
	oi := &entity.OrderItem{}
	query := "SELECT * FROM order_items WHERE order_id=$1"
	err := oit.db.Get(oi, query, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order item not found")
		}
		return nil, err
	}
	return oi, nil
}

// GetOrderItemByReferenceNumber retrieves an order item from the database by fulfillment_reference_number
func (oit *OrderItems) GetOrderItemByReferenceNumber(referenceNumber string) (*entity.OrderItem, error) {
	oi := &entity.OrderItem{}
	query := "SELECT * FROM order_items WHERE fulfillment_reference_number=$1"
	err := oit.db.Get(oi, query, referenceNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order item not found")
		}
		return nil, err
	}
	return oi, nil
}

// UpdateOrderItem updates an order item in the database
func (oit *OrderItems) UpdateOrderItem(oi *entity.OrderItem) error {
	query := "UPDATE Order_Items SET order_id=$1, product_id=$2, price=$4, cogs=$5, fulfillment_reference_number=$6, fulfillment_status=$7 WHERE id=$8"
	_, err := oit.db.Exec(query, oi.OrderID, oi.ProductID, oi.Price, oi.COGS, oi.FulfillmentReferenceNumber, oi.FulfillmentStatus, oi.ID)
	if err != nil {
		return err
	}
	return nil
}
