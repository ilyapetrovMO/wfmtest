package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := &dataJSON{
		"status": "available",
	}

	err := app.writeJSON(w, http.StatusOK, data)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}
