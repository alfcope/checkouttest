package errors

import (
	"fmt"
	"github.com/alfcope/checkout/model"
)

type ProductNotFound struct {
	Code model.ProductCode
}

type PromotionNotFound struct {
	Code model.PromotionType
}

type PromotionInvalid struct {
	Code string
	Msg  string
}

type BasketNotFound struct {
	Id string
}

type PrimaryKeyError struct {
	Id string
}

func NewProductNotFound(code model.ProductCode) *ProductNotFound {
	return &ProductNotFound{Code: code}
}

func NewPromotionNotFound(code model.PromotionType) *PromotionNotFound {
	return &PromotionNotFound{
		Code: code,
	}
}

func NewPromotionInvalid(code, message string) *PromotionInvalid {
	return &PromotionInvalid{
		Code: code,
		Msg:  message,
	}
}

func NewBasketNotFound(id string) *BasketNotFound {
	return &BasketNotFound{Id: id}
}

func NewPrimaryKeyError(id string) *PrimaryKeyError {
	return &PrimaryKeyError{Id: id}
}

//TODO: localization for error messages
func (p *ProductNotFound) Error() string {
	return fmt.Sprintf("Product %v not found", p.Code)
}

func (p *PromotionNotFound) Error() string {
	return fmt.Sprintf("Promotion %v not found", p.Code)
}

func (b *BasketNotFound) Error() string {
	return fmt.Sprintf("Basket %v not found", b.Id)
}

func (p *PromotionInvalid) Error() string {
	return fmt.Sprintf("Promotion %v invalid: %v", p.Code, p.Msg)
}

func (p *PrimaryKeyError) Error() string {
	return fmt.Sprintf("Primary key already exists: %v", p.Id)
}
