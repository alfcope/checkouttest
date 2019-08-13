package model

import (
	"github.com/google/uuid"
	"testing"
)

var testCases = []struct {
	itemBasket     []*Item   //Items in the basket
	promo          Promotion //Promotion to apply
	itWithoutPromo int       //Number of items out of the promotion
	itWithPromo    int       //Number of items eligible by the promotion
	free           int       //Number of free items if applicable, -1 otherwise
}{
	// ----- BULK PROMOTION TESTS ------
	{ //Edge case: basket without items
		[]*Item{},
		NewBulkPromotion(map[ProductCode]int{"P2": 1000}, 2),
		0,
		0,
		-1,
	}, {
		[]*Item{NewItem("P1", 1000), NewItem("P2", 1200), NewItem("P3", 1500)},
		NewBulkPromotion(map[ProductCode]int{"P2": 1000}, 2),
		3,
		0,
		-1,
	}, {
		[]*Item{NewItem("P1", 1000), NewItem("P1", 1000), NewItem("P1", 1000)},
		NewBulkPromotion(map[ProductCode]int{"P1": 850}, 3),
		0,
		3,
		-1,
	}, {
		[]*Item{NewItem("P2", 1200), NewItem("P2", 1200), NewItem("P2", 1200)},
		NewBulkPromotion(map[ProductCode]int{"P2": 1000}, 2),
		1,
		2,
		-1,
	}, {
		[]*Item{NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500),
			NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500),
			NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500),
			NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500)},
		NewBulkPromotion(map[ProductCode]int{"P3": 1100}, 4),
		3,
		12,
		-1,
	}, {
		[]*Item{NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500),
			NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500),
			NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500),
			NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500), NewItem("P3", 1500)},
		NewBulkPromotion(map[ProductCode]int{"P3": 1100}, 4),
		3,
		12,
		-1,
	},
	// ----- FREE ITEMS PROMOTION TESTS ------
	{ //Edge case: basket without items
		[]*Item{},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 2, Free: 1}}}),
		0,
		0,
		0,
	}, { //Different products
		[]*Item{NewItem("P2", 1500), NewItem("P2", 1500)},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 2, Free: 1}}}),
		2,
		0,
		0,
	}, {
		[]*Item{NewItem("P2", 800), NewItem("P2", 800)},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P2": {{Buy: 2, Free: 1}}}),
		0,
		2,
		1,
	}, {
		[]*Item{NewItem("P1", 1100), NewItem("P1", 1100)},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 2, Free: 1}}}),
		0,
		2,
		1,
	}, {
		[]*Item{NewItem("P1", 1500), NewItem("P1", 1500), NewItem("P1", 1500)},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 2, Free: 1}}}),
		1,
		2,
		1,
	}, { //Multiple promotions for same product
		[]*Item{NewItem("P1", 1500), NewItem("P1", 1500), NewItem("P1", 1500),
			NewItem("P1", 1500), NewItem("P1", 1500), NewItem("P1", 1500), NewItem("P1", 1500)},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 5, Free: 2}, {Buy: 2, Free: 1}}}),
		0,
		7,
		3,
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

		if tc.free >= 0 {
			freeCounter := 0
			if basket.itemsInPromotion[tc.promo.GetType()] != nil {
				for i := 0; i < len(*basket.itemsInPromotion[tc.promo.GetType()]); i++ {
					if (*basket.itemsInPromotion[tc.promo.GetType()])[i].finalPrice == 0 {
						freeCounter++
					}
				}
			}
			if tc.free != freeCounter {
				t.Errorf("got %v free items, want %v", freeCounter, tc.free)
			}
		}
	}
}
