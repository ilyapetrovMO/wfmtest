package main

import (
	"net/http"

	"github.com/ilyapetrovMO/WFMTtest/internal/token"
)

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	tokstr, err := app.getTokenFromHeader(&r.Header)
	if err != nil {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}

	claims, err := token.ParseJWT(tokstr)
	if err == token.ErrTokenNotValid {
		app.unauthorizedResponse(w, r, "token not valid")
		return
	}
	if err != nil {
		app.internalServerErrorResponse(w, r)
		return
	}

	orderbody := &struct {
		Product_id int
		Amount     int
	}{}
	err = app.readJSON(r, orderbody)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
	}

	order, err := app.models.Orders.CreateOrder(r.Context(), int64(claims.User_id), int64(orderbody.Product_id), int64(orderbody.Amount))
	if err != nil {
		app.internalServerErrorResponse(w, r)
		return
	}

	app.writeJSON(w, http.StatusOK, &dataJSON{"order": order})
}
