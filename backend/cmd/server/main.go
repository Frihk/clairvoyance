package main

import (
	"log"
	"net/http"

	"clairvoyance/internal/app"
	"clairvoyance/internal/config"
)

func main() {
	handler := app.MustHandler()
	if err := http.ListenAndServe(":"+config.App.Port, handler); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
