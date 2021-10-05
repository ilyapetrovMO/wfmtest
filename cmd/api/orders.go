package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
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

	order, err := app.models.Orders.CreateOrder(r.Context(), int64(claims.UserId), int64(orderbody.Product_id), int64(orderbody.Amount), time.Now())
	if err != nil {
		app.internalServerErrorResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"order": order})
	if err != nil {
		app.internalServerErrorResponse(w, r)
		return
	}
}

func (app *application) getOrdersForUser(w http.ResponseWriter, r *http.Request) {
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

	usrClaim, err := token.ParseJWT(tokstr)
	if err == token.ErrTokenNotValid {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}
	if err != nil {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}

	if usrIdParam != int64(usrClaim.UserId) {
		app.unauthorizedResponse(w, r, "unauthorized")
		return
	}

	orders, err := app.models.Orders.GetOrdersForUser(r.Context(), int64(usrClaim.UserId))
	if err == db.ErrRecordNotFound {
		app.writeJSON(w, http.StatusOK, &dataJSON{"orders": [0]int{}})
		return
	}
	if err != nil {
		app.internalServerErrorResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"orders": orders})
	if err != nil {
		app.internalServerErrorResponse(w, r)
		return
	}
}

func (app *application) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	tokstr, err := app.getTokenFromHeader(&r.Header)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	usrClaim, err := token.ParseJWT(tokstr)
	if err == token.ErrTokenNotValid {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}
	if err != nil {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}

	if usrClaim.RoleId != db.ROLE_MANAGER {
		app.unauthorizedResponse(w, r, "unauthorized")
		return
	}

	orders, err := app.models.Orders.GetOrders(r.Context())
	if err == db.ErrRecordNotFound {
		app.writeJSON(w, http.StatusOK, &dataJSON{"orders": [0]int{}})
		return
	}
	if err != nil {
		app.internalServerErrorResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"orders": orders})
	if err != nil {
		app.internalServerErrorResponse(w, r)
		return
	}
}

func (app *application) CancelOrder(w http.ResponseWriter, r *http.Request) {
	ordIdParam, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, "no valid query parameters given")
		return
	}

	tokStr, err := app.getTokenFromHeader(&r.Header)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	usrClaim, err := token.ParseJWT(tokStr)
	if err == token.ErrTokenNotValid {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}
	if err != nil {
		app.unauthorizedResponse(w, r, err.Error())
		return
	}

	order, err := app.models.Orders.GetOrderById(r.Context(), ordIdParam)
	if err == db.ErrRecordNotFound {
		app.badRequestResponse(w, r, "no record with specified id")
		return
	}
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		app.internalServerErrorResponse(w, r)
		return
	}

	if usrClaim.UserId != int(order.UserId) {
		app.unauthorizedResponse(w, r, "unauthorized")
		return
	}

	err = app.models.Orders.CancellOrder(r.Context(), order.OrderId, time.Now())
	if err == db.ErrRecordNotFound {
		app.badRequestResponse(w, r, "no record with specified id")
	}
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		app.internalServerErrorResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"status": "ok"})
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		app.internalServerErrorResponse(w, r)
		return
	}
}
