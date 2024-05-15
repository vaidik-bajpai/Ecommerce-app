package main

import (
	"net/http"

	"github.com/vaidik-bajpai/ecommerce-api/internal/validator"
)

func (app *application) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	v := validator.NewValidator()

	userId := app.readQueryInt(qs, "user_id", v)
	productId := app.readQueryInt(qs, "product_id", v)

	if userId == -1 || productId == -1 || !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err := app.models.Carts.AddToCart(userId, productId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "The product is added to cart"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) removeItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("removeItemHandler endpoint"))
}

func (app *application) cartCheckoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("cartCheckoutHandler endpoint"))
}

func (app *application) instantBuyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("instantBuyHandler endpoint"))
}
