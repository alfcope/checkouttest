package model

import (
	"sync"
)

type Basket struct {
	Id    string
	lines map[ProductCode]Line

	rwMux sync.RWMutex
}

type Line struct {
	Product
	amount int
}

func NewBasket(id string) *Basket {
	return &Basket{
		Id:    id,
		lines: make(map[ProductCode]Line),
		rwMux: sync.RWMutex{},
	}
}

func (b *Basket) AddProduct(p Product) {
	b.rwMux.Lock()
	defer b.rwMux.Unlock()

	if l, ok := b.lines[p.Code]; ok {
		l.amount++
		return
	}

	b.lines[p.Code] = Line{
		Product: p,
		amount:  1,
	}
}

func (b *Basket) CalculatePrice(offers []Promotion) float64 {
	var productInOffer = make(map[ProductCode]*[]int)
	var price = 0

	b.rwMux.Lock()
	defer b.rwMux.Unlock()

	if offers != nil && len(offers) > 0 {
		for _, p := range offers {
			p.Resolve(b.lines, productInOffer)
		}
	}

	for pcode, line := range b.lines {
		inOfferCounter := 0
		if inOffer, ok := productInOffer[pcode]; ok {
			if inOffer != nil && len(*inOffer) > 0 {
				inOfferCounter = len(*inOffer)
				for _, offerPrice := range *inOffer {
					price += offerPrice
				}
			}
		}

		price += (line.amount - inOfferCounter) * line.Price
	}

	return float64(price / 100)
}
