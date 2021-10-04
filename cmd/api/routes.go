package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()

	mux.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	mux.HandlerFunc(http.MethodPost, "/auth", app.authenticationHandler)
	mux.HandlerFunc(http.MethodGet, "/products", app.getAllProductsHandler)
	mux.Handler(http.MethodPost, "/orders", app.userOnly(http.HandlerFunc(app.createOrderHandler)))
	mux.Handler(http.MethodPost, "/products", app.managerOnly(http.HandlerFunc(app.createProductHandler)))

	return mux
}
