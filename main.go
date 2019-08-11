package main

import (
	"github.com/alfcope/checkout/config"
	"github.com/alfcope/checkout/pkg/logging"
	"github.com/alfcope/checkout/server"
)

func main() {
	configuration, err := config.LoadConfiguration("./config", "configuration")

	api, err := server.NewCheckoutApi(configuration)
	if err != nil {
		logging.Logger.Error("Shutting down. Error initialing api: ", err.Error())
		return
	}

	api.RunServer(configuration.Server.Port)
}
