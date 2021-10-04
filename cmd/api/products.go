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

func (app *application) getAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := app.models.Products.GetProducts(r.Context())
	if err != nil {
		app.internalServerErrorResponse(w, r)
	}

	if len(products) == 0 {
		app.writeJSON(w, http.StatusOK, &dataJSON{"products": []struct{}{}})
	}

	dtoarr := []*Product{}
	for _, pr := range products {
		dtoarr = append(dtoarr, &Product{Name: pr.Name, Description: pr.Description})
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"products": dtoarr})
	if err != nil {
		app.internalServerErrorResponse(w, r)
	}
}
