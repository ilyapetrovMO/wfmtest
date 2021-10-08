package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
)

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InStorage   int    `json:"in_storage"`
}

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	prjs := &Product{}
	err := app.readJSON(w, r, prjs)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "unable to parse json")
		return
	}

	v := validator.New()

	product := &db.Product{
		Name:        prjs.Name,
		Description: prjs.Description,
		InStorage:   int64(prjs.InStorage),
	}

	db.ValidateProduct(v, product)

	if !v.Valid() {
		app.badRequestResponse(w, r, v.Errors())
		return
	}

	err = app.models.Products.Create(product)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			app.alreadyExistsResponse(w, r, "product with this name already exists")
			return
		}

		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"product": product})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := app.models.Products.GetAll()
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"products": products})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	input := &struct {
		ProductId   int64     `json:"product_id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		InStorage   int64     `json:"in_storage"`
	}{}

	err := app.readJSON(w, r, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product, err := app.models.Products.GetById(input.ProductId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			app.badRequestResponse(w, r, "no product with specified id")
			return
		}

		app.internalServerErrorResponse(w, r, err)
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.InStorage = input.InStorage

	v := validator.New()

	db.ValidateProduct(v, product)

	if !v.Valid() {
		app.badRequestResponse(w, r, v.Errors())
		return
	}

	err = app.models.Products.Update(product)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"product": product})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}
