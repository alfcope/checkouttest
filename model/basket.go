package model

import (
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

	total float32
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
