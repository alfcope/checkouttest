package api

import (
	"github.com/alfcope/checkouttest/datasource"
	"github.com/alfcope/checkouttest/model"
	"github.com/google/uuid"
)

type checkoutService struct {
	ds datasource.Datasource
}

type CheckoutService interface {
	CreateBasket() (string, error)
	AddItem(string, model.ProductCode) error
	DeleteBasket(string)
}

func NewCheckoutService(ds datasource.Datasource) CheckoutService {
	return &checkoutService{
		ds: ds,
	}
}

func (c *checkoutService) CreateBasket() (string, error) {
	id := uuid.New().String()

	basket := model.NewBasket(id)

	err := c.ds.AddBasket(basket)
	//TODO: hash collision control!!
	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *checkoutService) AddItem(basketId string, productCode model.ProductCode) error {

	_, err := c.ds.GetProduct(productCode)
	if err != nil {
		return err
	}

	return c.ds.AddItemToBasket(basketId, productCode)
}

func (c *checkoutService) DeleteBasket(id string) {
	c.ds.DeleteBasket(id)
}
