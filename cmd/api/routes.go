package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	mux.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	mux.HandlerFunc(http.MethodPost, "/user/login", app.loginUserHandler)
	mux.HandlerFunc(http.MethodPost, "/user/register", app.registerUserHandler)

	mux.Handler(http.MethodPost, "/cart/item", app.userOnly(http.HandlerFunc(app.addToCartHandler)))
	mux.Handler(http.MethodDelete, "/cart/item", app.userOnly(http.HandlerFunc(app.removeFromCartHandler)))
	mux.Handler(http.MethodPost, "/cart/checkout", app.userOnly(http.HandlerFunc(app.checkoutCartHandler)))
	mux.Handler(http.MethodGet, "/cart", app.userOnly(http.HandlerFunc(app.getUserCartAndItems)))
	mux.Handler(http.MethodGet, "/carts", app.managerOnly(http.HandlerFunc(app.getAllCartsAndItems)))

	mux.HandlerFunc(http.MethodGet, "/product", app.getAllProductsHandler)
	mux.Handler(http.MethodPost, "/product", app.managerOnly(http.HandlerFunc(app.createProductHandler)))
	mux.Handler(http.MethodPut, "/product", app.managerOnly(http.HandlerFunc(app.updateProductHandler)))

	mux.Handler(http.MethodGet, "/order", app.managerOnly(http.HandlerFunc(app.getAllOrders)))
	mux.Handler(http.MethodGet, "/order/user/:id", app.userOnly(http.HandlerFunc(app.getOrdersForUserId)))
	mux.Handler(http.MethodDelete, "/order/:id", app.userOnly(http.HandlerFunc(app.deleteOrder)))

	return app.getJWT(mux)
}
