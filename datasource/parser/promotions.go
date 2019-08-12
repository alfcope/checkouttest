package parser

import (
	"fmt"
	"github.com/alfcope/checkouttest/errors"
	"github.com/alfcope/checkouttest/model"
)

func ParsePromotion(nodes map[string]interface{}) (model.Promotion, error) {

	//fmt.Printf("%T - %v \n", nodes, nodes)

	switch nodes["code"].(string) {
	case "BULK":
		return parseBulkPromotion(nodes)

	case "ONE_FREE":
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

func parseOneFreePromotion(nodes map[string]interface{}) (*model.OneFreePromotion, error) {
	var items map[model.ProductCode]int

	rawItems := nodes["items"].([]interface{})
	items = make(map[model.ProductCode]int, len(rawItems))

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
		if _, ok := item["amount"].(float64); !ok {
			fmt.Printf("Invalid amount: %v", item["amount"])
			continue
		}

		items[model.ProductCode(item["product"].(string))] = int(item["amount"].(float64))
	}

	if len(items) == 0 {
		return nil, errors.NewPromotionInvalid(nodes["code"].(string), "empty items list")
	}

	return model.NewOneFreeProduction(items), nil
}
