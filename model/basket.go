package model

import (
	"fmt"
	"strconv"
	"sync"
)

type Item struct {
	code       ProductCode
	finalPrice int
}

type Basket struct {
	Id               string
	items            map[ProductCode][]Item
	itemsInPromotion map[PromotionType]*[]Item

	total float64
	RWMux sync.RWMutex
}

func NewItem(pcode ProductCode, price int) *Item {
	return &Item{
		code:       pcode,
		finalPrice: price,
	}
}

func NewBasket(id string) *Basket {
	return &Basket{
		Id:               id,
		items:            make(map[ProductCode][]Item),
		itemsInPromotion: make(map[PromotionType]*[]Item),
		total:            -1,
		RWMux:            sync.RWMutex{},
	}
}

func (b *Basket) AddItem(newItem *Item) {
	b.RWMux.Lock()
	defer b.RWMux.Unlock()

	b.items[newItem.code] = append(b.items[newItem.code], *newItem)
}

func (b *Basket) CalculatePrice() (float64, error) {
	total := 0

	b.RWMux.Lock()
	defer b.RWMux.Unlock()

	if b.total > -1 {
		return b.total, nil
	}

	//Items out of any promotion
	for _, items := range b.items {
		if items != nil {
			total += items[0].finalPrice * len(items)
		}
	}

	//Items in promotions
	for _, items := range b.itemsInPromotion {
		for _, item := range *items {
			total += item.finalPrice
		}
	}

	totalPrice, err := strconv.ParseFloat(fmt.Sprintf("%.2f", total),64)
	if err != nil {
		return 0, err
	}

	return totalPrice, nil
}
