package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := &dataJSON{
		"status": "available",
	}

	app.writeJSON(w, http.StatusOK, data)
}
