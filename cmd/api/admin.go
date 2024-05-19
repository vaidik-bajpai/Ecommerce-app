package main

import (
	"errors"
	"net/http"

	"github.com/vaidik-bajpai/ecommerce-api/internal/data"
	"github.com/vaidik-bajpai/ecommerce-api/internal/validator"
)

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name   string `json:"name"`
		Price  uint64 `json:"price"`
		Rating uint8  `json:"rating"`
		Image  string `json:"image"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product := &data.Product{
		Name:   input.Name,
		Price:  input.Price,
		Rating: input.Rating,
		Image:  &input.Image,
	}

	v := validator.NewValidator()

	if data.ValidateProduct(v, product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Products.AddProduct(product)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, product, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundErrorResponse(w, r)
		return
	}

	err = app.models.Products.RemoveProduct(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "product successfullt deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
