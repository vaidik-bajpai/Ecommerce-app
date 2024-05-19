package main

import (
	"errors"
	"math"
	"net/http"

	"github.com/vaidik-bajpai/ecommerce-api/internal/data"
	"github.com/vaidik-bajpai/ecommerce-api/internal/validator"
)

func (app *application) searchProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundErrorResponse(w, r)
		return
	}

	product, err := app.models.Products.Get(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) searchProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string
		Price uint
		data.Filters
	}

	qs := r.URL.Query()

	v := validator.NewValidator()

	input.Name = app.readString(qs, "product_name", "", v)
	input.Price = uint(app.readInt(qs, "price", math.MaxInt32, v))

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Sort = app.readString(qs, "sort", "id", v)
	input.Filters.SortSafeList = []string{"id", "name", "price", "-id", "rating", "-name", "-price", "-rating"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	products, metadata, err := app.models.Products.GetAll(input.Name, int(input.Price), input.Filters)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
