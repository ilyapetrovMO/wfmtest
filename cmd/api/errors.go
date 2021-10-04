package main

import "net/http"

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, msg string) {
	err := app.writeJSON(w, status, &dataJSON{"error": msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusUnauthorized, "unauthorized")
}

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusInternalServerError, "internal server error")
}
