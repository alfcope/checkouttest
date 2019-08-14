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
	AddProduct(string, model.ProductCode) error
	GetBasketPrice(string) (float64, error)
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

func (c *checkoutService) AddProduct(bId string, pCode model.ProductCode) error {

	p, err := c.ds.GetProduct(pCode)
	if err != nil {
		return err
	}

	basket, err := c.ds.GetBasket(bId)
	if err != nil {
		return err
	}

	basket.AddProduct(p)

	return nil
}

func (c *checkoutService) GetBasketPrice(id string) (float64, error) {

	return 0, nil
}

func (c *checkoutService) DeleteBasket(id string) {
	c.ds.DeleteBasket(id)
}
