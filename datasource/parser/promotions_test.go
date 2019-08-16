package parser

import (
	"github.com/alfcope/checkouttest/errors"
	"github.com/alfcope/checkouttest/model"
	"reflect"
	"testing"
)

var promotionsParsersCases = []struct {
	nodes     map[string]interface{}
	promotion model.Promotion
	err       error
}{
	// ---- BULK PROMOTIONS CASES
	{ // Empty promotion
		map[string]interface{}{},
		model.NewBulkPromotion(map[model.ProductCode][]model.BulkOfferRule{}),
		errors.NewPromotionNotFound(""),
	}, { // Promotion without promos
		map[string]interface{}{"code": "BULK", "promos": []interface{}{}},
		nil,
		errors.NewPromotionInvalid("BULK", "empty items list"),
	}, { // Promotion with a wrong product code
		map[string]interface{}{"code": "BULK", "promos": []interface{}{
			map[string]interface{}{"product": "", "rules": []interface{}{map[string]interface{}{"buy": float64(3), "price": float64(1000)},
				map[string]interface{}{"buy": float64(-5), "price": float64(850)},},
			},
			map[string]interface{}{"product": "PR2", "rules": []interface{}{map[string]interface{}{"buy": float64(3), "price": float64(500)},},},
		}},
		model.NewBulkPromotion(map[model.ProductCode][]model.BulkOfferRule{
			"PR2": {{3, 500}},
		}, ),
		nil,
	}, { // Promotion with a wrong buy value
		map[string]interface{}{"code": "BULK", "promos": []interface{}{
			map[string]interface{}{"product": "PR1", "rules": []interface{}{map[string]interface{}{"buy": float64(3), "price": float64(1000)},
				map[string]interface{}{"buy": float64(-5), "price": float64(850)},},
			},
			map[string]interface{}{"product": "PR2", "rules": []interface{}{map[string]interface{}{"buy": float64(3), "price": float64(500)},},},
		}},
		model.NewBulkPromotion(map[model.ProductCode][]model.BulkOfferRule{
			"PR1": {{3, 1000}},
			"PR2": {{3, 500}},
		}, ),
		nil,
	}, { // Promotion with a wrong price value
		map[string]interface{}{"code": "BULK", "promos": []interface{}{
			map[string]interface{}{"product": "PR1", "rules": []interface{}{map[string]interface{}{"buy": float64(3), "price": float64(-1000)},
				map[string]interface{}{"buy": float64(5), "price": float64(850)},},
			},
			map[string]interface{}{"product": "PR2", "rules": []interface{}{map[string]interface{}{"buy": float64(3), "price": float64(500)},},},
		}},
		model.NewBulkPromotion(map[model.ProductCode][]model.BulkOfferRule{
			"PR1": {{5, 850}},
			"PR2": {{3, 500}},
		}, ),
		nil,
	}, { // Promotion with a promotion without rules
		map[string]interface{}{"code": "BULK", "promos": []interface{}{
			map[string]interface{}{"product": "PR1", "rules": []interface{}{},
			},
			map[string]interface{}{"product": "PR2", "rules": []interface{}{map[string]interface{}{"buy": float64(2), "price": float64(600)},},},
		}},
		model.NewBulkPromotion(map[model.ProductCode][]model.BulkOfferRule{
			"PR2": {{2, 600}},
		}, ),
		nil,
	}, { // Correct promotion
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
			if err.Error() != pc.err.Error() {
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
