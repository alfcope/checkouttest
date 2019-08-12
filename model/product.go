package model

type ProductCode string

type Product struct {
	Code  ProductCode `json:"code"`
	Name  string      `json:"name"`
	Price int         `json:"price"`
}

func NewProduct(code ProductCode, name string, price int) *Product {
	return &Product{
		Code:  code,
		Name:  name,
		Price: price,
	}
}
