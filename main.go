package main

import (
	"github.com/alfcope/checkouttest/config"
	"github.com/alfcope/checkouttest/pkg/logging"
	"github.com/alfcope/checkouttest/server"
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
