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
