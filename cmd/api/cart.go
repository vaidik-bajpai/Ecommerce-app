package main

import "net/http"

func (app *application) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("addToCartHandler endpoint"))
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
