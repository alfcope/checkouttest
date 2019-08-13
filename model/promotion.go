package model

import (
	"log"
)

type PromotionType string

type Promotion interface {
	GetType() PromotionType
	ApplyTo(*Basket)
	Resolve(map[ProductCode][]Item) map[PromotionType][]Item
}

type BulkPromotion struct {
	//A map in case different bulk promotions are defined for different products
	//Key: ProductCode
	//Value: different possible conditions by product, for example => 3 - $19 | 5 - $15
	offers  map[ProductCode][]BulkOfferRule
}

type BulkOfferRule struct {
	Buy  int
	Price int
}

func NewBulkPromotion(offers map[ProductCode][]BulkOfferRule) *BulkPromotion {
	return &BulkPromotion{
		offers: offers,
	}
}

func (b BulkPromotion) GetType() PromotionType {
	return "BULK"
}

func (b BulkPromotion) ApplyTo(basket *Basket) {
	//basket.RWMux.Lock()
	//defer basket.RWMux.Unlock()
	//
	//log.Printf("--- Bulk Promotion ---")
	//for p, price := range b.offers {
	//	log.Println("\tProduct: ", p)
	//
	//	if bitems, ok := basket.items[p]; ok {
	//		log.Printf("\tFound %v in the basket\n", len(bitems))
	//
	//		promotions := len(bitems) / b.amount
	//		log.Printf("\tEnough items for %v promotions\n", promotions)
	//		if promotions > 0 {
	//			elements := promotions * b.amount
	//			log.Printf("\t%v elements to move\n", elements)
	//			basket.items[p] = basket.items[p][elements:]
	//			log.Printf("\t%v without promotion: %v\n", p, len(basket.items[p]))
	//
	//			if basket.itemsInPromotion[b.GetType()] != nil {
	//				log.Printf("\t%v in promotion: %v\n", p, len(*basket.itemsInPromotion[b.GetType()]))
	//			}
	//
	//			if basket.itemsInPromotion[b.GetType()] == nil {
	//				items := make([]Item, 0, elements)
	//				basket.itemsInPromotion[b.GetType()] = &items
	//				log.Printf("\tSlice created! %v", elements)
	//			}
	//			for elements > 0 {
	//				tmp := append(*basket.itemsInPromotion[b.GetType()], *NewItem(p, price))
	//				log.Printf("\ttmp length: %v\n", len(tmp))
	//
	//				basket.itemsInPromotion[b.GetType()] = &tmp
	//				elements--
	//				log.Printf("\t%v elements remaining\n", elements)
	//			}
	//			log.Printf("\t%v in promotion %v: %v\n", p, b.GetType(), len(*basket.itemsInPromotion[b.GetType()]))
	//		}
	//	}
	//}
}

func (b BulkPromotion) Resolve(basket map[ProductCode][]Item) map[PromotionType][]Item {
	return make(map[PromotionType][]Item)
}


type FreeItemsPromotion struct {
	//A map in case different bulk promotions are defined for different products
	//Key: ProductCode
	//Value: slice with potentially different combinations of buy X get Y free
	offers map[ProductCode][]FreeItemsOfferRule
}

type FreeItemsOfferRule struct {
	Buy  int
	Free int
}

func NewFreeItemsPromotion(offers map[ProductCode][]FreeItemsOfferRule) *FreeItemsPromotion {
	return &FreeItemsPromotion{offers: offers}
}

func (f FreeItemsPromotion) GetType() PromotionType {
	return "FREE_ITEMS"
}

func (f FreeItemsPromotion) ApplyTo(basket *Basket) {
	//basket.RWMux.Lock()
	//defer basket.RWMux.Unlock()
	//
	//log.Printf("--- Free Items Promotion ---")
	//for pcode, conditions := range f.items {
	//	log.Println("\tProduct: ", pcode)
	//
	//	if bitems, ok := basket.items[pcode]; ok {
	//		log.Printf("\tFound %v in the basket\n", len(bitems))
	//
	//		for _, condition := range conditions {
	//			numPromoEligible := len(basket.items[pcode]) / condition.Buy
	//
	//			log.Printf("\tEnough items for %v promotions\n", numPromoEligible)
	//			if numPromoEligible > 0 {
	//				productPrice := basket.items[pcode][0].finalPrice
	//
	//				elements := numPromoEligible * condition.Buy
	//				log.Printf("\t%v elements to move\n", elements)
	//				basket.items[pcode] = basket.items[pcode][elements:]
	//				log.Printf("\t%v without promotion: %v\n", pcode, len(basket.items[pcode]))
	//
	//				if basket.itemsInPromotion[f.GetType()] != nil {
	//					log.Printf("\t%v in promotion: %v\n", pcode, len(*basket.itemsInPromotion[f.GetType()]))
	//				}
	//
	//				if basket.itemsInPromotion[f.GetType()] == nil {
	//					items := make([]Item, 0, elements)
	//					basket.itemsInPromotion[f.GetType()] = &items
	//					log.Printf("\tSlice created! %v", elements)
	//				}
	//
	//				for i := 0; i < elements; i++ {
	//					price := productPrice
	//					if i < condition.Free {
	//						price = 0
	//					}
	//
	//					tmp := append(*basket.itemsInPromotion[f.GetType()], *NewItem(pcode, price))
	//
	//					log.Printf("\ttmp length: %v\n", len(tmp))
	//
	//					basket.itemsInPromotion[f.GetType()] = &tmp
	//					log.Printf("\t%v elements remaining\n", elements-i-1)
	//				}
	//			}
	//		}
	//	}
	//}
}

func (f FreeItemsPromotion) Resolve(basket map[ProductCode][]Item) map[PromotionType][]Item {

	itemsInPromotions := make(map[PromotionType][]Item)

	log.Printf("--- Free Items Promotion ---")
	for pcode, conditions := range f.offers {
		log.Println("\tProduct: ", pcode)

		if bitems, ok := basket[pcode]; ok {
			log.Printf("\tFound %v in the basket\n", len(bitems))

			inPromotionCounter := 0
			for _, condition := range conditions {
				numPromoEligible := (len(basket[pcode]) - inPromotionCounter) / condition.Buy

				log.Printf("\tEnough items for %v promotions\n", numPromoEligible)
				if numPromoEligible > 0 {
					productPrice := basket[pcode][0].finalPrice

					elements := numPromoEligible * condition.Buy
					inPromotionCounter += elements

					log.Printf("\t%v elements to move\n", elements)
					//basket.items[pcode] = basket.items[pcode][elements:]
					//log.Printf("\t%v without promotion: %v\n", pcode, len(basket.items[pcode]))

					if itemsInPromotions[f.GetType()] == nil {
						items := make([]Item, 0, elements)
						itemsInPromotions[f.GetType()] = items
						log.Printf("\tSlice created! %v", elements)
					}

					for i := 0; i < elements; i++ {
						price := productPrice
						if i < condition.Free {
							price = 0
						}

						tmp := append(itemsInPromotions[f.GetType()], *NewItem(pcode, price))

						log.Printf("\ttmp length: %v\n", len(tmp))

						itemsInPromotions[f.GetType()] = tmp
						log.Printf("\t%v elements remaining\n", elements-i-1)
					}
				}
			}
		}
	}

	return itemsInPromotions
}
