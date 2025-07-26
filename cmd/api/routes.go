package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/v1/healthcheck", app.healthcheckHandler)

	mux.Get("/v1/ws", app.chatHandler)

	return mux
}
