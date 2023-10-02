package psql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gitlab.com/lokalpay-dev/digital-goods/internal/entity"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/slog"
)

type Customer struct {
	db *sqlx.DB
}

func NewCustomer(db *sqlx.DB) *Customer {
	return &Customer{
		db: db,
	}
}

// CreateCustomer inserts a new customer into the database
func (c *Customer) CreateCustomer(email, passwordHash string) (int, error) {
	var customerID int
	query := "INSERT INTO customers (email, password, created_at) VALUES ($1, $2, NOW()) RETURNING id"
	row := c.db.QueryRow(query, email, passwordHash)
	err := row.Scan(&customerID)
	if err != nil || customerID == 0 {
		return 0, err
	}

	return customerID, nil
}

// GetCustomerByEmail retrieves a customer from the database by Email
func (c *Customer) GetCustomerByEmail(email string) (entity.Customer, error) {
	cst := entity.Customer{}
	query := "SELECT * FROM customers WHERE email=$1"
	err := c.db.Get(&cst, query, email)
	if err != nil && err != sql.ErrNoRows {
		slog.Errorw("unexpected error", "stack_trace", err.Error())
		return cst, err
	}
	return cst, nil
}

// GetCustomerByEmail retrieves a customer from the database by Email
func (c *Customer) GetCustomerByID(id int) (entity.Customer, error) {
	cst := entity.Customer{}
	query := "SELECT * FROM customers WHERE id=$1"
	err := c.db.Get(&cst, query, id)
	if err != nil && err != sql.ErrNoRows {
		slog.Errorw("unexpected error", "stack_trace", err.Error())
		return cst, err
	}
	return cst, nil
}

// UpdateCustomer updates a customer in the database
func (c *Customer) UpdateCustomer(cst *entity.Customer) error {
	query := "UPDATE customers SET name=$1, email=$3, phone_number=$4 WHERE id=$5"
	_, err := c.db.Exec(query, cst.Name, cst.Email, cst.PhoneNumber, cst.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCustomer deletes a customer from the database by ID
func (c *Customer) DeleteCustomer(id int) error {
	query := "DELETE FROM customers WHERE ID=$1"
	_, err := c.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
