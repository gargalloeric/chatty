package main

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelop{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
		},
	}

	if err := app.writeJSON(w, http.StatusOK, data, nil); err != nil {
		app.logError(r, err)
		app.serverErrorResponse(w, r, err)
		return
	}
}
