package model

import (
	"testing"
)

var promotionCases = []struct {
	basketLines    map[ProductCode]Line //Items in the basket
	promo          Promotion            //Promotion to apply
	itWithoutPromo int                  //Number of items out of the promotion
	itWithPromo    int                  //Number of items eligible by the promotion
	free           int                  //Number of free items if applicable, -1 otherwise
}{
	// ----- BULK PROMOTION TESTS ------
	{ //Edge case: basket without items
		make(map[ProductCode]Line),
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P2": {{Buy: 2, Price: 1000,}}}),
		0,
		0,
		-1,
	}, {
		map[ProductCode]Line{"P1": {Product{"P1", "aaaa", 1000,}, 1,},
			"P2": {Product{"P2", "bbbb", 1200,}, 1,},
			"P3": {Product{"P3", "cccc", 1500,}, 1,}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P2": {{Buy: 2, Price: 1000,}}}),
		3,
		0,
		-1,
	}, { //Exact amount of items
		map[ProductCode]Line{"P1": {Product{"P1", "aaaa", 1000,}, 3,}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{Buy: 3, Price: 850,}}}),
		0,
		3,
		-1,
	}, {
		map[ProductCode]Line{"P2": {Product{"P2", "bbbb", 1200,}, 3,}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P2": {{Buy: 2, Price: 1000,}}}),
		1,
		2,
		-1,
	}, { //Exact amount of same items matching two different rules
		map[ProductCode]Line{"P1": {Product{"P1", "aaaa", 1000,}, 7,}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{Buy: 5, Price: 650,}, {Buy: 2, Price: 850,}}}),
		0,
		7,
		-1,
	}, {
		map[ProductCode]Line{"P3": {Product{"P3", "cccc", 1500,}, 15,}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P3": {{Buy: 4, Price: 1100,}}}),
		3,
		12,
		-1,
	}, { //Exact amount of two different items matching two different rules
		map[ProductCode]Line{"P1": {Product{"P1", "aaaa", 1500,}, 3,},
			"P2": {Product{"P2", "bbbb", 1200,}, 3,}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{Buy: 3, Price: 1300,},},
			"P2": {{Buy: 3, Price: 1000,},}}),
		0,
		6,
		-1,
	},
	// ----- FREE ITEMS PROMOTION TESTS ------
	{ //Edge case: basket without items
		make(map[ProductCode]Line),
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P2": {{Buy: 2, Free: 1,}}}),
		0,
		0,
		0,
	}, {
		map[ProductCode]Line{"P1": {Product{"P1", "aaaa", 1000,}, 1,},
			"P2": {Product{"P2", "bbbb", 1200,}, 1,},
			"P3": {Product{"P3", "cccc", 1500,}, 1,}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P2": {{Buy: 2, Free: 1,}}}),
		3,
		0,
		0,
	}, { //Exact amount of items
		map[ProductCode]Line{"P1": {Product{"P1", "aaaa", 1000,}, 3,}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 3, Free: 1,}}}),
		0,
		3,
		1,
	}, {
		map[ProductCode]Line{"P2": {Product{"P2", "bbbb", 1200,}, 3,}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P2": {{Buy: 2, Free: 1,}}}),
		1,
		2,
		1,
	}, { //Exact amount of same items matching two different rules
		map[ProductCode]Line{"P1": {Product{"P1", "aaaa", 1000,}, 7,}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 5, Free: 3,}, {Buy: 2, Free: 1,}}}),
		0,
		7,
		4,
	}, {
		map[ProductCode]Line{"P3": {Product{"P3", "cccc", 1500,}, 15,}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P3": {{Buy: 4, Free: 1,}}}),
		3,
		12,
		3,
	}, { //Exact amount of two different items matching two different rules
		map[ProductCode]Line{"P1": {Product{"P1", "aaaa", 1500,}, 3,},
			"P2": {Product{"P2", "bbbb", 1200,}, 3,}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 3, Free: 1,},},
			"P2": {{Buy: 3, Free: 1,},}}),
		0,
		6,
		2,
	},
}

func TestPromotions(t *testing.T) {
	for _, tc := range promotionCases {
		inOffer := make(map[ProductCode]*[]int)

		tc.promo.Resolve(tc.basketLines, inOffer)

		cInOffer := 0
		cOutOffer := 0
		freeCounter := 0
		for pCode, line := range tc.basketLines {
			if items, ok := inOffer[pCode]; ok {
				cInOffer += len(*items)
				cOutOffer += line.amount - len(*items)

				if tc.free >= 0 {
					for _, price := range *items {
						if price == 0 {
							freeCounter++
						}
					}
				}
			} else {
				cOutOffer += line.amount
			}
		}

		if cInOffer != tc.itWithPromo {
			t.Errorf("got %v items with promo, want %v", cInOffer, tc.itWithPromo)
		}
		if cOutOffer != tc.itWithoutPromo {
			t.Errorf("got %v items without promo, want %v", cOutOffer, tc.itWithoutPromo)
		}
		if tc.free >= 0 && freeCounter != tc.free {
			t.Errorf("got %v items free, want %v", freeCounter, tc.free)
		}

	}
}
