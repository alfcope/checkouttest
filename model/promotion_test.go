package model

import (
	"github.com/google/uuid"
	"testing"
)

var testCases = []struct {
	itemBasket     []*Item
	promo          Promotion
	itWithoutPromo int
	itWithPromo    int
}{
	{
		[]*Item{NewItem("P1", 1000), NewItem("P2", 1200), NewItem("P3", 1500)},
		NewBulkPromotion(map[ProductCode]int{"P2": 1000}, 2),
		3,
		0,
	}, {
		[]*Item{NewItem("P1", 1000), NewItem("P1", 1000), NewItem("P1", 1000)},
		NewBulkPromotion(map[ProductCode]int{"P1": 850}, 3),
		0,
		3,
	},{
		[]*Item{NewItem("P2", 1200), NewItem("P2", 1200), NewItem("P2", 1200)},
		NewBulkPromotion(map[ProductCode]int{"P2": 1000}, 2),
		1,
		2,
	}, {
		[]*Item{NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500),
			NewItem("P3", 1500),NewItem("P3", 1500),NewItem("P3", 1500),NewItem("P3", 1500),
			NewItem("P3", 1500),NewItem("P3", 1500),NewItem("P3", 1500),NewItem("P3", 1500),
			NewItem("P3", 1500),NewItem("P3", 1500),NewItem("P3", 1500),NewItem("P3", 1500)},
		NewBulkPromotion(map[ProductCode]int{"P3": 1100}, 4),
		3,
		12,
	},
}

func TestPromotions(t *testing.T) {
	for _, tc := range testCases {
		basket := NewBasket(uuid.New().String())
		for _, it := range tc.itemBasket {
			basket.AddItem(it)
		}

		tc.promo.ApplyTo(basket)

		var counter int
		for _, v := range basket.items {
			counter += len(v)
		}
		if counter != tc.itWithoutPromo {
			t.Errorf("got %v items without promo, want %v", counter, tc.itWithoutPromo)
		}

		if basket.itemsInPromotion[tc.promo.GetType()] == nil {
			if tc.itWithPromo != 0 {
				t.Errorf("got 0 items with promo, want %v", tc.itWithPromo)
			}
		} else if len(*basket.itemsInPromotion[tc.promo.GetType()]) != tc.itWithPromo {
			t.Errorf("got %v items with promo, want %v", len(*basket.itemsInPromotion[tc.promo.GetType()]), tc.itWithPromo)
		}
	}
}