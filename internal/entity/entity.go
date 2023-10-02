package entity

import "time"

type Customer struct {
	ID          int        `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Email       string     `db:"email" json:"email"`
	Password    string     `db:"password" json:"-"`
	PhoneNumber *string    `db:"phone_number" json:"phone_number"`
	IsVerified  bool       `db:"is_verified" json:"is_verified"`
	CreatedAt   *time.Time `db:"created_at" json:"-"`
	UpdatedAt   *time.Time `db:"updated_at" json:"-"`
}

type Category struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

type Product struct {
	ID          int        `db:"id" json:"id"`
	ProductName string     `db:"product_name" json:"product_name"`
	CategoryID  int        `db:"category_id" json:"-"`
	Description *string    `db:"description" json:"description"`
	COGS        float64    `db:"cogs" json:"cogs"`
	PartnerFee  float64    `db:"partner_fee" json:"partner_fee"`
	CreatedAt   *time.Time `db:"created_at" json:"-"`
	UpdatedAt   *time.Time `db:"updated_at" json:"-"`
}

type Order struct {
	ID         int        `db:"id" json:"id"`
	CustomerID int        `db:"customer_id" json:"customer_id"`
	OrderDate  time.Time  `db:"order_date" json:"order_date"`
	TotalPrice float64    `db:"total_price" json:"total_price"`
	Status     string     `db:"status" json:"status"` // PENDING,SUCCESS,FAILED
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	ID                         int        `db:"id" json:"id"`
	OrderID                    int        `db:"order_id" json:"order_id"`
	ProductID                  int        `db:"product_id" json:"product_id"`
	ClientNumber               string     `db:"client_number" json:"client_number"`
	Price                      float64    `db:"price" json:"price"`
	COGS                       float64    `db:"cogs" json:"cogs"`
	FulfillmentStatus          string     `db:"fulfillment_status"`           // PENDING,SUCCESS,FAILED
	FulfillmentReferenceNumber string     `db:"fulfillment_reference_number"` // transaction_id between lokalpay and ppob provider
	SerialCode                 *string    `db:"serial_code"`
	CreatedAt                  *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                  *time.Time `db:"updated_at" json:"updated_at"`
}

type Payment struct {
	ID                int        `db:"id" json:"id"`
	OrderID           int        `db:"order_id" json:"order_id"`
	PaymentDate       *time.Time `db:"payment_date" json:"payment_date"`
	PaymentMethod     string     `db:"payment_method" json:"payment_method"` // VIRTUAL_ACCOUNT - CREDIT_CARD - E_WALLET
	PaymentExpiryTime *time.Time `db:"payment_expiry_time" json:"payment_expiry_time"`
	PaymentAmount     float64    `db:"payment_amount" json:"payment_amount"`
	AdminFee          float64    `db:"admin_fee" json:"admin_fee"`
	Status            string     `db:"status" json:"status"`                     // PENDING,SUCCESS,FAILED
	ReferenceNumber   string     `db:"reference_number" json:"reference_number"` // transaction_id between lokalpay and payment provider
	AccountNumber     string     `db:"account_number" json:"account_number"`
	CreatedAt         *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at" json:"updated_at"`
}
