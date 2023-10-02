package service

import (
	"gitlab.com/lokalpay-dev/digital-goods/internal"
	"gitlab.com/lokalpay-dev/digital-goods/internal/constant"
	"gitlab.com/lokalpay-dev/digital-goods/internal/entity"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/helper"
)

type Product struct {
	productRepo internal.ProductRepositoryItf
}

func NewProduct(
	productRepo internal.ProductRepositoryItf,
) *Product {
	return &Product{
		productRepo: productRepo,
	}
}

func (p *Product) GetProducts() ([]entity.Product, error) {
	// select from product using constant category_id
	return p.productRepo.GetProducts()
}

func (p *Product) OperatorSelection(inputNumber string) string {
	var opt string
	if helper.Contains(constant.TselPrefix, inputNumber) {
		opt = constant.PLTelkomsepInstitutionCode
	} else if helper.Contains(constant.XlPrefix, inputNumber) {
		opt = constant.PLXLInstitutionCode
	} else if helper.Contains(constant.IM3Prefix, inputNumber) {
		opt = constant.PLIndosatInstitutionCode
	} else if helper.Contains(constant.TriPrefix, inputNumber) {
		opt = constant.PLTriInstitutionCode
	} else {
		opt = ""
	}

	return opt
}
