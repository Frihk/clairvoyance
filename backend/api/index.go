package handler

import (
	"log"
	"net/http"
	"sync"

	"clairvoyance/internal/app"
)

var (
	handler http.Handler
	once    sync.Once
	initErr error
)

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		handler, initErr = app.NewHandler()
		if initErr != nil {
			log.Printf("failed to initialize backend handler: %v", initErr)
		}
	})

	if initErr != nil {
		http.Error(w, initErr.Error(), http.StatusInternalServerError)
		return
	}

	handler.ServeHTTP(w, r)
}
