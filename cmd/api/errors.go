package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message any) {
	data := envelop{"error": message}

	if err := app.writeJSON(w, statusCode, data, nil); err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered an error and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
