package datasource

import (
	"encoding/json"
	"github.com/alfcope/checkouttest/config"
	"github.com/alfcope/checkouttest/datasource/parser"
	"github.com/alfcope/checkouttest/errors"
	"github.com/alfcope/checkouttest/model"
	"io/ioutil"
	"log"
	"sync"
)

type Datasource struct {
	// products and promotions do not need mutex as they do not
	// change its state. Just once at startup
	products   map[model.ProductCode]model.Product
	promotions map[model.PromotionType]model.Promotion

	baskets    map[string]model.Basket
	basketsMux sync.RWMutex
}

func InitDatasource(config config.DataConfig) (*Datasource, error) {
	ds := Datasource{
		products:   make(map[model.ProductCode]model.Product),
		promotions: make(map[model.PromotionType]model.Promotion),
		baskets:    make(map[string]model.Basket),
		basketsMux: sync.RWMutex{},
	}

	err := ds.loadProducts(config.Products)
	if err != nil {
		return nil, err
	}

	err = ds.loadPromotions(config.Promotions)
	if err != nil {
		return nil, err
	}

	return &ds, nil
}

func (d *Datasource) GetProduct(code model.ProductCode) (model.Product, error) {
	if product, ok := d.products[code]; ok {
		return product, nil
	}

	return *new(model.Product), errors.NewProductNotFound(code)
}

func (d *Datasource) GetPromotion(code model.PromotionType) (model.Promotion, error) {
	if promotion, ok := d.promotions[code]; ok {
		return promotion, nil
	}

	return *new(model.Promotion), errors.NewPromotionNotFound(code)
}

func (d *Datasource) GetPromotions() map[model.PromotionType]model.Promotion {
	return d.promotions
}

func (d *Datasource) GetBasket(id string) (*model.Basket, error) {
	if basket, ok := d.baskets[id]; ok {
		return &basket, nil
	}

	return new(model.Basket), errors.NewBasketNotFound(id)
}

func (d *Datasource) AddBasket(basket *model.Basket) error {
	d.basketsMux.Lock()
	defer d.basketsMux.Unlock()

	if _, ok := d.baskets[basket.Id]; !ok {
		d.baskets[basket.Id] = *basket
		return nil
	}

	return errors.NewPrimaryKeyError(basket.Id)
}

func (d *Datasource) AddItemToBasket(basketId string, code model.ProductCode) error {

	if _, ok := d.products[code]; !ok {
		return errors.NewProductNotFound(code)
	}

	log.Println("Adding - Baskets size: ", len(d.baskets))

	if basket, ok := d.baskets[basketId]; ok {
		basket.AddItem(model.NewItem(code, d.products[code].Price))
		return nil
	}

	return errors.NewBasketNotFound(basketId)
}

func (d *Datasource) DeleteBasket(basketId string) {
	d.basketsMux.Lock()
	defer d.basketsMux.Unlock()

	delete(d.baskets, basketId)
	log.Println("Deleting - Baskets size: ", len(d.baskets))
}

func (d *Datasource) loadProducts(filePath string) error {
	var products []model.Product

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &products)
	if err != nil {
		return err
	}

	for _, p := range products {
		d.products[p.Code] = p
	}
	return nil
}

func (d *Datasource) loadPromotions(filePath string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var nodes []map[string]interface{}
	err = json.Unmarshal(file, &nodes)
	if err != nil {
		return err
	}

	for _, promotionNode := range nodes {
		promotion, err := parser.ParsePromotion(promotionNode)
		if err != nil {
			if _, ok := err.(*errors.PromotionNotFound); !ok {
				return err
			}
			continue
		}

		d.promotions[promotion.GetType()] = promotion
	}

	return nil
}
