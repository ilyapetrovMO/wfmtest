package main

import "net/http"

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, msg string) {
	err := app.writeJSON(w, status, &dataJSON{"error": msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, msg string) {
	app.errorResponse(w, r, http.StatusUnauthorized, msg)
}

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, msg string) {
	app.errorResponse(w, r, http.StatusBadRequest, msg)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, "Method Not Allowed")
}
