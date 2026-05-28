package main

import (
	"clairvoyance/internal/api/router"
	"clairvoyance/internal/config"
	"clairvoyance/internal/utils"
	"log"
)

func main() {
	r := router.SetupRouter()

	utils.InitLogger(config.App.AppEnv)

	if err := r.Run(":9000"); err != nil {
		log.Fatal("Failed to start")
	}
}
