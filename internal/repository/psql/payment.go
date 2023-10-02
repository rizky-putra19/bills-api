package psql

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"gitlab.com/lokalpay-dev/digital-goods/internal/entity"
)

type Payment struct {
	db *sqlx.DB
}

func NewPayment(db *sqlx.DB) *Payment {
	return &Payment{
		db: db,
	}
}

// CreatePayment inserts a new payment into the database
func (pmt *Payment) CreatePayment(p *entity.Payment) error {
	query := `INSERT INTO payments (
		order_id, payment_date, payment_method, payment_amount,
		admin_fee, payment_expiry_time, account_number, reference_number, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	row := pmt.db.QueryRow(query,
		p.OrderID, p.PaymentDate, p.PaymentMethod,
		p.PaymentAmount, p.AdminFee, p.PaymentExpiryTime,
		p.AccountNumber, p.ReferenceNumber, p.Status,
	)
	err := row.Scan(&p.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetPayment retrieves a payment from the database by ID
func (pmt *Payment) GetPayment(id int) (*entity.Payment, error) {
	p := &entity.Payment{}
	query := "SELECT * FROM payments WHERE id=$1"
	err := pmt.db.Get(p, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found") // payment not found
		}
		return nil, err
	}
	return p, nil
}

// GetPaymentByReferenceNumber retrieves a payment from the database by reference_number
func (pmt *Payment) GetPaymentByReferenceNumber(referenceNumber string) (*entity.Payment, error) {
	p := &entity.Payment{}
	query := "SELECT * FROM payments WHERE reference_number=$1"
	err := pmt.db.Get(p, query, referenceNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found") // payment not found
		}
		return nil, err
	}
	return p, nil
}

// UpdatePayment updates a payment in the database
func (pmt *Payment) UpdatePayment(p *entity.Payment) error {
	query := `
	UPDATE 
    	payments 
	SET 
	    order_id=$1, payment_date=$2, payment_method=$3,
	    payment_amount=$4, admin_fee=$5, status=$6, account_number=$7, updated_at=now()
	WHERE id=$8`
	_, err := pmt.db.Exec(query, p.OrderID, p.PaymentDate, p.PaymentMethod, p.PaymentAmount, p.AdminFee, p.Status, p.AccountNumber, p.ID)
	if err != nil {
		return err
	}
	return nil
}
