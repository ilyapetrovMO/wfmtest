package main

import (
	"net/http"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/ilyapetrovMO/WFMTtest/internal/token"
	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
)

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	input := &struct {
		ProductId int `json:"product_id"`
		Amount    int `json:"amount"`
		UserId    int `json:"user_id"`
	}{}

	err := app.readJSON(w, r, input)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

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
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if claims.UserId != input.UserId {
		app.unauthorizedResponse(w, r, "unauthorized")
		return
	}

	order := &db.Order{
		ProductId: int64(input.ProductId),
		Amount:    int64(input.Amount),
		UserId:    int64(input.UserId),
	}

	v := validator.New()

	db.ValidateOrder(v, order)

	if !v.Valid() {
		app.badRequestResponse(w, r, v.Errors())
		return
	}

	err = app.models.Orders.Create(order)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"order": order})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) getOrdersForUserId(w http.ResponseWriter, r *http.Request) {
	usrIdParam, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, "no valid query parameters given")
		return
	}

	tokstr, err := app.getTokenFromHeader(&r.Header)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	claims, err := token.ParseJWT(tokstr)
	if err == token.ErrTokenNotValid {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if claims.RoleId != db.ROLE_MANAGER && usrIdParam != int64(claims.UserId) {
		app.unauthorizedResponse(w, r, "unauthorized")
		return
	}

	orders, err := app.models.Orders.GetWithUserId(int64(claims.UserId))
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"orders": orders})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAllOrders(w http.ResponseWriter, r *http.Request) {
	tokstr, err := app.getTokenFromHeader(&r.Header)
	if err != nil {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}

	claims, err := token.ParseJWT(tokstr)
	if err != nil {
		app.unauthorizedResponse(w, r, "invalid token")
		return
	}

	if claims.RoleId != db.ROLE_MANAGER {
		app.unauthorizedResponse(w, r, "unauthorized")
		return
	}

	orders, err := app.models.Orders.GetAll()
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"orders": orders})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteOrder(w http.ResponseWriter, r *http.Request) {
	orderIdParam, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	tokStr, err := app.getTokenFromHeader(&r.Header)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	claims, err := token.ParseJWT(tokStr)
	if err == token.ErrTokenNotValid {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}
	if err != nil {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}

	order, err := app.models.Orders.GetById(orderIdParam)
	if err == db.ErrRecordNotFound {
		app.badRequestResponse(w, r, "no record with specified id")
		return
	}
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if claims.UserId != int(order.UserId) {
		app.unauthorizedResponse(w, r, "unauthorized")
		return
	}

	err = app.models.Orders.Delete(order)
	if err == db.ErrRecordNotFound {
		app.badRequestResponse(w, r, "no record with specified id")
	}
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"status": "ok"})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}
