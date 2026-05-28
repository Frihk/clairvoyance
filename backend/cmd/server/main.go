package main

import (
	"clairvoyance/internal/api/router"
	"log"
)

func main() {
	r := router.SetupRouter()

	if err := r.Run(":9000"); err != nil {
		log.Fatal("Failed to start")
	}
}
