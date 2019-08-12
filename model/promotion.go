package model

import "log"

type PromotionType string

type Promotion interface {
	GetType() PromotionType
	ApplyTo(basket *Basket)
}

type BulkPromotion struct {
	//A map in case different bulk promotions are defined for different products
	//Key: ProductCode
	//Value: product price
	items  map[ProductCode]int
	amount int
}

func NewBulkPromotion(items map[ProductCode]int, amount int) *BulkPromotion {
	return &BulkPromotion{
		items:  items,
		amount: amount,
	}
}

func (b BulkPromotion) GetType() PromotionType {
	return "BULK"
}

func (b BulkPromotion) ApplyTo(basket *Basket) {
	log.Printf("--- Bulk Promotion ---")
	for p, price := range b.items {
		log.Println("\tProduct: ", p)

		if bitems, ok := basket.items[p]; ok {
			log.Printf("\tFound %v in the basket\n", len(bitems))

			promotions := len(bitems) / b.amount
			log.Printf("\tEnough items for %v promotions\n", promotions)
			if promotions > 0 {
				elements := promotions * b.amount
				log.Printf("\t%v elements to move\n", elements)
				basket.items[p] = basket.items[p][elements:]
				log.Printf("\t%v without promotion: %v\n", p, len(basket.items[p]))

				if basket.itemsInPromotion[b.GetType()] != nil {
					log.Printf("\t%v in promotion: %v\n", p, len(*basket.itemsInPromotion[b.GetType()]))
				}

				if basket.itemsInPromotion[b.GetType()] == nil {
					items := make([]Item, 0, elements)
					basket.itemsInPromotion[b.GetType()] = &items
					log.Printf("\tSlice created! %v", elements)
				}
				for elements > 0 {
					tmp := append(*basket.itemsInPromotion[b.GetType()], *NewItem(p, price))
					log.Printf("\ttmp length: %v\n", len(tmp))

					basket.itemsInPromotion[b.GetType()] = &tmp
					elements--
					log.Printf("\t%v elements remaining\n", elements)
				}
				log.Printf("\t%v in promotion %v: %v\n", p, b.GetType(), len(*basket.itemsInPromotion[b.GetType()]))
			}
		}
	}
}

type FreeItemsPromotion struct {
	//A map in case different bulk promotions are defined for different products
	//Key: ProductCode
	//Value: slice with potentially different combinations of buy X get Y free
	items map[ProductCode][]FreeItemsPromoConditions
}

type FreeItemsPromoConditions struct {
	Buy  int
	Free int
}

func NewFreeItemsPromotion(items map[ProductCode][]FreeItemsPromoConditions) *FreeItemsPromotion {
	return &FreeItemsPromotion{items: items}
}

func (f FreeItemsPromotion) GetType() PromotionType {
	return "FREE_ITEMS"
}

func (f FreeItemsPromotion) ApplyTo(basket *Basket) {
	log.Printf("--- Free Items Promotion ---")
	for p, conditions := range f.items {
		log.Println("\tProduct: ", p)

		if bitems, ok := basket.items[p]; ok {
			log.Printf("\tFound %v in the basket\n", len(bitems))

			for _, condition := range conditions {
				numPromoEligible := len(basket.items[p]) / condition.Buy

				log.Printf("\tEnough items for %v promotions\n", numPromoEligible)
				if numPromoEligible > 0 {
					productPrice := basket.items[p][0].finalPrice

					elements := numPromoEligible * condition.Buy
					log.Printf("\t%v elements to move\n", elements)
					basket.items[p] = basket.items[p][elements:]
					log.Printf("\t%v without promotion: %v\n", p, len(basket.items[p]))

					if basket.itemsInPromotion[f.GetType()] != nil {
						log.Printf("\t%v in promotion: %v\n", p, len(*basket.itemsInPromotion[f.GetType()]))
					}

					if basket.itemsInPromotion[f.GetType()] == nil {
						items := make([]Item, 0, elements)
						basket.itemsInPromotion[f.GetType()] = &items
						log.Printf("\tSlice created! %v", elements)
					}

					for i := 0; i < elements; i++ {
						price := productPrice
						if i < condition.Free {
							price = 0
						}

						tmp := append(*basket.itemsInPromotion[f.GetType()], *NewItem(p, price))

						log.Printf("\ttmp length: %v\n", len(tmp))

						basket.itemsInPromotion[f.GetType()] = &tmp
						log.Printf("\t%v elements remaining\n", elements-i-1)
					}
				}
			}
		}
	}
}
