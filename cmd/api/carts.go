package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
)

func (app *application) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	input := &struct {
		UserId    int `json:"user_id"`
		ProductId int `json:"product_id"`
		Amount    int `json:"amount"`
	}{}

	err := app.readJSON(w, r, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	claims, ok := r.Context().Value(userClaimsKey{}).(*userClaims)
	if !ok || input.UserId != claims.userId {
		app.unauthorizedResponse(w, r, "unauthorized")
		return
	}

	cart, err := app.models.Carts.GetByUserId(int64(claims.userId))
	if err != nil {
		app.internalServerErrorResponse(w, r, errors.New("expected to find cart, got none"))
		return
	}

	newItem := &db.CartItem{
		CartId:    cart.CartId,
		ProductId: int64(input.ProductId),
		Amount:    int64(input.Amount),
	}

	v := validator.New()

	db.ValidateCartItem(v, newItem)

	if !v.Valid() {
		app.badRequestResponse(w, r, v.Errors())
		return
	}

	product, err := app.models.Products.GetById(int64(input.ProductId))
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			app.badRequestResponse(w, r, "product does not exist")
			return
		}
		app.internalServerErrorResponse(w, r, err)
		return
	}

	exists := false

	existingItem, err := app.models.CartItem.GetByProductId(product.ProductId)
	if err != nil {
		if !errors.Is(err, db.ErrRecordNotFound) {
			app.internalServerErrorResponse(w, r, err)
			return
		}

	} else {
		exists = true
	}

	if exists {
		if existingItem.Amount+newItem.Amount > product.InStorage {
			app.badRequestResponse(w, r, fmt.Sprintf("not enough product in storage, %d available", product.InStorage))
			return
		}

		existingItem.Amount += newItem.Amount
		err := app.models.CartItem.Update(existingItem)
		if err != nil {
			app.internalServerErrorResponse(w, r, err)
			return
		}

		err = app.writeJSON(w, http.StatusCreated, &dataJSON{"item": existingItem})
		if err != nil {
			app.internalServerErrorResponse(w, r, err)
			return
		}

		return
	}

	if newItem.Amount > product.InStorage {
		app.badRequestResponse(w, r, fmt.Sprintf("not enough product in storage, %d available", product.InStorage))
		return
	}

	err = app.models.CartItem.Create(newItem)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, &dataJSON{"item": newItem})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	input := &struct {
		CartItemId int `json:"cart_item_id"`
		UserId     int `json:"user_id"`
		Amount     int `json:"amount"`
	}{}

	err := app.readJSON(w, r, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	claims, ok := r.Context().Value(userClaimsKey{}).(*userClaims)
	if !ok || input.UserId != claims.userId {
		app.unauthorizedResponse(w, r, "unauthorized")
		return
	}

	existingItem, _, err := app.models.CartItem.GetById(int64(input.CartItemId))
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			app.badRequestResponse(w, r, "no cart item with specified cart_item_id")
			return
		}
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if input.Amount <= 0 {
		app.badRequestResponse(w, r, "amount: must be a positive non-zero value")
		return
	}

	existingItem.Amount = int64(input.Amount)

	product, err := app.models.Products.GetById(existingItem.ProductId)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if existingItem.Amount > product.InStorage {
		app.badRequestResponse(w, r, fmt.Sprintf("not enough product in storage, %d available", product.InStorage))
		return
	}

	err = app.models.CartItem.Update(existingItem)
	if err != nil {
		app.internalServerErrorResponse(w, r, errors.New("could not update item"))
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"cart_item": existingItem})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) removeFromCartHandler(w http.ResponseWriter, r *http.Request) {
	itemIdParam, err := app.readCartItemIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	userIdParam, err := app.readUserIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	claims, ok := r.Context().Value(userClaimsKey{}).(*userClaims)
	if !ok || int(userIdParam) != claims.userId {
		app.unauthorizedResponse(w, r, "unathorized")
		return
	}

	item, userId, err := app.models.CartItem.GetById(itemIdParam)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			app.badRequestResponse(w, r, "item does not exist")
			return
		}

		app.internalServerErrorResponse(w, r, err)
		return
	}

	if userId != int(userIdParam) {
		app.unauthorizedResponse(w, r, "unathorized")
		return
	}

	err = app.models.CartItem.Delete(item)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			app.badRequestResponse(w, r, "no such cart item exists")
			return
		}
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"status": "cart item removed"})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

type Cart struct {
	db.Cart
	Items []*db.CartItem
}

func (app *application) getAllCartsAndItemsHandler(w http.ResponseWriter, r *http.Request) {
	carts, err := app.models.Carts.GetAll()
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	arr := []*Cart{}

	for _, c := range carts {
		newCart := &Cart{}
		newCart.CartId = c.CartId
		newCart.UserId = c.UserId
		newCart.CreatedAt = c.CreatedAt

		items, err := app.models.CartItem.GetAllByUserId(c.UserId)
		if err != nil {
			app.internalServerErrorResponse(w, r, err)
			return
		}

		newCart.Items = items

		arr = append(arr, newCart)
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"carts": arr})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) getUserCartAndItemsHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := app.readUserIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	claims, ok := r.Context().Value(userClaimsKey{}).(*userClaims)
	if !ok || int(userId) != claims.userId {
		app.unauthorizedResponse(w, r, "unauthorized")
		return
	}

	cart, err := app.models.Carts.GetByUserId(userId)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	newCart := &Cart{}
	newCart.CartId = cart.CartId
	newCart.UserId = cart.UserId
	newCart.CreatedAt = cart.CreatedAt

	items, err := app.models.CartItem.GetAllByUserId(userId)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	newCart.Items = items

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"cart": newCart})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) checkoutCartHandler(w http.ResponseWriter, r *http.Request) {
	cartId, err := app.readCartIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cart, err := app.models.Carts.GetById(cartId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			app.badRequestResponse(w, r, "no cart with specified id")
			return
		}

		app.internalServerErrorResponse(w, r, err)
		return
	}

	claims, ok := r.Context().Value(userClaimsKey{}).(*userClaims)
	if !ok || int(cart.UserId) != claims.userId {
		app.unauthorizedResponse(w, r, "unathorized")
		return
	}

	cartItems, err := app.models.CartItem.GetAllByUserId(int64(claims.userId))
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
	}

	if len(cartItems) == 0 {
		app.badRequestResponse(w, r, "cart is empty")
		return
	}

	orders := []*db.Order{}
	for _, ci := range cartItems {
		prd, err := app.models.Products.GetById(ci.ProductId)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				app.badRequestResponse(w, r, fmt.Sprintf("product with id %d not found", ci.ProductId))
				return
			}
			app.internalServerErrorResponse(w, r, err)
			return
		}

		order, err := app.models.Orders.CreateTx(ci, prd, int64(claims.userId))
		if err != nil {
			app.internalServerErrorResponse(w, r, err)
			return
		}

		orders = append(orders, order)
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"orders": orders})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}
