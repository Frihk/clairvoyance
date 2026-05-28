package main

import (
	"clairvoyance/internal/api/router"
	"clairvoyance/internal/config"
	"clairvoyance/internal/utils"
	"log"
)

func main() {
	config.Load()
	utils.InitLogger(config.App.AppEnv)

	r := router.SetupRouter()
	if err := r.Run(":" + config.App.Port); err != nil {
		log.Fatal("Failed to start:", err)
	}
}
