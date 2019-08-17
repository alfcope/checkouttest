package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alfcope/checkouttest/api/responses"
	"io/ioutil"
	"net/http"
	"time"
)

type CheckoutClient struct {
	serverUrl  string
	apiVersion string
	httpClient *http.Client
}

func NewCheckoutClient(serverUrl, version string) *CheckoutClient {
	return &CheckoutClient{
		serverUrl:  serverUrl,
		apiVersion: version,
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *CheckoutClient) AddBasket() (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/api/%v/baskets/", c.serverUrl, c.apiVersion), nil)
	if err != nil {
		return "", fmt.Errorf("there was an error creating http request: %v", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("%v", resp.Status)
	}

	if resp.Body != nil {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error fetching response body: %v", err)
		}

		nb := responses.NewBasketResponse{}
		err = json.Unmarshal(responseBody, &nb)
		if err != nil {
			return "", fmt.Errorf("error fetching response body: %v", err)
		}

		return nb.Id, nil
	}

	return "", errors.New("empty response")
}
