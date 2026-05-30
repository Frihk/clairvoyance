package handler

import (
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
	})

	if initErr != nil {
		http.Error(w, initErr.Error(), http.StatusInternalServerError)
		return
	}

	handler.ServeHTTP(w, r)
}
