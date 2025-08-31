package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFoundResponse)
	mux.MethodNotAllowed(app.methodNotAllowedResponse)

	mux.Get("/v1/healthcheck", app.healthcheckHandler)

	mux.Route("/v1/rooms", func(r chi.Router) {
		r.Get("/", app.listRoomHandler)
		r.Post("/", app.createRoomHandler)
		r.Get("/{roomID}", app.connectToRoomHandler)
	})

	return mux
}
