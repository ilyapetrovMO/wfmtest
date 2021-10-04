package main

import "net/http"

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	prod := &Product{}
	err := app.readJSON(r, prod)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "unable to parse json")
		return
	}

	newprod, err := app.models.Products.CreateProduct(r.Context(), prod.Name, prod.Description)
	if err != nil {
		app.internalServerErrorResponse(w, r)
		return
	}

	datajs := &dataJSON{
		"product": &Product{
			Name:        newprod.Name,
			Description: newprod.Description,
		},
	}
	err = app.writeJSON(w, http.StatusOK, datajs)
	if err != nil {
		app.internalServerErrorResponse(w, r)
		return
	}
}
