package parser

import (
	"github.com/alfcope/checkouttest/model"
	"reflect"
	"testing"
)

var promotionsParsersCases = []struct {
	nodes     map[string]interface{}
	promotion model.Promotion
	err       error
}{
	{ // Correct promotion
		map[string]interface{}{"code": "BULK", "promos": []interface{}{
			map[string]interface{}{"product": "PR1", "rules": []interface{}{map[string]interface{}{"buy": float64(3), "price": float64(1000)},
				map[string]interface{}{"buy": float64(5), "price": float64(850)},},
			},
			map[string]interface{}{"product": "PR2", "rules": []interface{}{map[string]interface{}{"buy": float64(3), "price": float64(500)},},},
		}},
		model.NewBulkPromotion(map[model.ProductCode][]model.BulkOfferRule{
			"PR1": {{3, 1000}, {5, 850}},
			"PR2": {{3, 500}},
		}, ),
		nil,
	},
}

func TestBasketPrices(t *testing.T) {
	for _, pc := range promotionsParsersCases {

		promotion, err := ParsePromotion(pc.nodes)
		if err != nil {
			if pc.err == nil {
				t.Errorf("Unexpected error: %v", err.Error())
			}
			if err != pc.err {
				t.Errorf("Got error: %v, wanted: %v", err.Error(), pc.err.Error())
			}
			return
		}

		if pc.err != nil {
			t.Errorf("Did not get expected error: %v", pc.err.Error())
		}

		if !reflect.DeepEqual(promotion, pc.promotion) {
			t.Errorf("Got promotion %v, wanted %v", promotion, pc.promotion)
		}
	}
}
