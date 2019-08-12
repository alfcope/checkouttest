package parser

import (
	"fmt"
	"github.com/alfcope/checkouttest/errors"
	"github.com/alfcope/checkouttest/model"
)

func ParsePromotion(nodes map[string]interface{}) (model.Promotion, error) {

	switch nodes["code"].(string) {
	case "BULK":
		return parseBulkPromotion(nodes)

	case "FREE_ITEMS":
		return parseOneFreePromotion(nodes)

	default:
		return nil, errors.NewPromotionNotFound(model.PromotionType(nodes["code"].(string)))
	}
}

func parseBulkPromotion(nodes map[string]interface{}) (*model.BulkPromotion, error) {
	var items map[model.ProductCode]int
	amount := nodes["amount"]

	rawItems := nodes["items"].([]interface{})
	items = make(map[model.ProductCode]int, len(rawItems))

	for _, rawItem := range rawItems {
		if _, ok := rawItem.(map[string]interface{}); !ok {
			fmt.Printf("Invalid map: %v", rawItem)
			continue
		}
		item := rawItem.(map[string]interface{})
		if _, ok := item["product"].(string); !ok {
			fmt.Printf("Invalid product code: %v %T\n", item["product"], item["product"])
			continue
		}
		if _, ok := item["price"].(float64); !ok {
			fmt.Printf("Invalid price: %v\n", item["price"])
			continue
		}

		items[model.ProductCode(item["product"].(string))] = int(item["price"].(float64))
	}

	if len(items) == 0 {
		return nil, errors.NewPromotionInvalid(nodes["code"].(string), "empty items list")
	}

	return model.NewBulkPromotion(items, int(amount.(float64))), nil
}

func parseOneFreePromotion(nodes map[string]interface{}) (*model.FreeItemsPromotion, error) {
	items := make(map[model.ProductCode][]model.FreeItemsPromoConditions)

	rawItems := nodes["items"].([]interface{})

	for _, rawItem := range rawItems {
		if _, ok := rawItem.(map[string]interface{}); !ok {
			fmt.Printf("Invalid map: %v", rawItem)
			continue
		}
		item := rawItem.(map[string]interface{})

		if _, ok := item["product"].(string); !ok {
			fmt.Printf("Invalid product code: %v", item["product"])
			continue
		}
		if _, ok := item["buy"].(float64); !ok {
			fmt.Printf("Invalid amount: %v", item["buy"])
			continue
		}
		if _, ok := item["free"].(float64); !ok {
			fmt.Printf("Invalid amount: %v", item["free"])
			continue
		}

		promoConditions := model.FreeItemsPromoConditions{
			Buy:  int(item["buy"].(float64)),
			Free: int(item["free"].(float64)),
		}

		if promosProducto, ok := items[model.ProductCode(item["product"].(string))]; ok {
			items[model.ProductCode(item["product"].(string))] = append(promosProducto, promoConditions)
		} else {
			items[model.ProductCode(item["product"].(string))] = []model.FreeItemsPromoConditions{promoConditions}
		}
	}

	if len(items) == 0 {
		return nil, errors.NewPromotionInvalid(nodes["code"].(string), "empty items list")
	}

	return model.NewFreeItemsPromotion(items), nil
}
