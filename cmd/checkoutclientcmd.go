package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alfcope/checkouttest/cli"
	"github.com/alfcope/checkouttest/model"
	"github.com/manifoldco/promptui"
	"io/ioutil"
)

type RequestType int

const (
	GoBack RequestType = iota
	AddBasket
	AddProduct
	GetPrice
	DeleteBasket
)

type Operation struct {
	requestType RequestType
	description string
}

type CheckoutCmd struct {
	operations   []Operation
	basketIds    [] string
	productCodes [] string

	client *cli.CheckoutClient
}

func NewCheckoutCmd(productsPath, serverAddress string, apiVersion int) *CheckoutCmd {
	operations := []Operation{{
		GoBack, "Exit",
	}, {
		AddBasket, "Add new basket",
	}, {
		AddProduct, "Add new product to a basket",
	}, {
		GetPrice, "Get a basket price",
	}, {
		DeleteBasket, "Delete a basket",
	},}

	cmd := CheckoutCmd{
		operations:   operations,
		basketIds:    []string{operations[0].description},
		productCodes: []string{operations[0].description},
		client:       cli.NewCheckoutClient(serverAddress, apiVersion),
	}

	err := cmd.loadProducts(fmt.Sprintf("%sproducts.json", productsPath))
	if err != nil {
		fmt.Printf("Error loading products: %v", err.Error())
		return nil
	}

	return &cmd
}

func main() {
	productsPath := flag.String("products", "./config/", "path to folder containing the available list of products file")
	serverAddress := flag.String("server", "http://localhost:7070", "server http address")
	apiVersion := flag.Int("version", 1, "api version to request")

	flag.Parse()

	cmd := NewCheckoutCmd(*productsPath, *serverAddress, *apiVersion)
	if cmd == nil {
		return
	}

	for {
		exit := cmd.showInitialScreen()

		if exit {
			break
		}
	}
}

func (c *CheckoutCmd) loadProducts(filePath string) error {
	var products []model.Product

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &products)
	if err != nil {
		return err
	}

	for _, p := range products {
		err := p.Validate()
		if err == nil {
			c.productCodes = append(c.productCodes, string(p.Code))
		}
	}

	return nil
}

func (c *CheckoutCmd) showInitialScreen() bool {

	prompt := promptui.Select{
		Label: "Select Option",
		Items: c.operations,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return false
	}

	switch i {
	case 0:
		return true
	case 1:
		id, err := c.client.AddBasket()
		if err != nil {
			fmt.Printf("Error adding basket: %v\n", err)
		}
		c.basketIds = append(c.basketIds, id)
		fmt.Printf("Basket %v added\n", id)

	case 2:
		c.showBasketsList(AddProduct)
	case 3:
		c.showBasketsList(GetPrice)
	case 4:
		c.showBasketsList(DeleteBasket)
	}

	return false
}

func (c *CheckoutCmd) showBasketsList(requestType RequestType) {
	prompt := promptui.Select{
		Label: "Select Basket",
		Items: c.basketIds,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if i == 0 {
		return
	}

	switch requestType {
	case GetPrice:
		c.showPrice(c.basketIds[i])

	case DeleteBasket:
		err = c.client.DeleteBasket(c.basketIds[i])
		if err != nil {
			fmt.Printf("Error deleting basket %v: %v", c.basketIds[i], err.Error())
		} else {
			fmt.Printf("Basket %v deleted!\n", c.basketIds[i])
			c.basketIds = remove(c.basketIds, i)
		}

	default:
		c.showProductLists(c.basketIds[i])
	}
}

func (c *CheckoutCmd) showProductLists(basketId string) {
	for {
		prompt := promptui.Select{
			Label: "Select Product",
			Items: c.productCodes,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if i == 0 {
			return
		}

		err = c.client.AddItem(basketId, c.productCodes[i])
		if err != nil {
			fmt.Printf("Error adding product: %v\n", err)
		}
		fmt.Printf("%v added to basket %v", c.productCodes[i], basketId)
	}
}

func (c *CheckoutCmd) showPrice(basketId string) {
	price, err := c.client.GetPrice(basketId)
	if err != nil {
		fmt.Printf("Error getting price: %v\n", err)
	}
	fmt.Printf("Basket %v price: %.2f\n", basketId, price)
}

func remove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
