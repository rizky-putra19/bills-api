package psql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gitlab.com/lokalpay-dev/digital-goods/internal/entity"
)

type Product struct {
	db *sqlx.DB
}

func NewProduct(db *sqlx.DB) *Product {
	return &Product{
		db: db,
	}
}

func (prd *Product) GetProduct(productID int) (*entity.Product, error) {
	p := &entity.Product{}
	query := "SELECT * FROM products WHERE id=$1"
	err := prd.db.Get(p, query, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // product not found
		}
		return nil, err
	}
	return p, nil
}

func (prd *Product) GetProducts() ([]entity.Product, error) {
	p := []entity.Product{}
	query := "SELECT * FROM products ORDER BY cogs ASC"
	err := prd.db.Select(&p, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // product not found
		}
		return nil, err
	}
	return p, nil
}
