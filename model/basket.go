package model

import (
	"log"
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
		total:            -1, //Price cache
		RWMux:            sync.RWMutex{},
	}
}

func (b *Basket) AddItem(newItem *Item) {
	b.RWMux.Lock()
	defer b.RWMux.Unlock()

	b.items[newItem.code] = append(b.items[newItem.code], *newItem)
	// Expire price cache
	b.total = -1
}

func (b *Basket) GetTotal() float64 {
	return b.total
}

func (b *Basket) CalculatePrice() (float64, error) {
	total := 0

	log.Print("Waiting for!")
	b.RWMux.Lock()
	defer b.RWMux.Unlock()
	log.Print("Lock acquired!")

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

	b.total = float64(total/100)

	return b.total, nil
}
