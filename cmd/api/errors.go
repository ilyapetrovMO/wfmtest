package main

import "net/http"

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	err := app.writeJSON(w, status, &dataJSON{"error": data})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, msg string) {
	w.Header().Set("XXX-Authorization", "Bearer")
	app.errorResponse(w, r, http.StatusUnauthorized, msg)
}

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err)
	app.errorResponse(w, r, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	app.errorResponse(w, r, http.StatusBadRequest, data)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, "Method Not Allowed")
}

func (app *application) alreadyExistsResponse(w http.ResponseWriter, r *http.Request, msg string) {
	app.errorResponse(w, r, http.StatusConflict, msg)
}
