package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/users/signup", app.createUserHandler)
	mux.HandleFunc("POST /v1/users/login", app.userLoginHandler)
	mux.HandleFunc("POST /v1/users/productview", app.viewProductHandler)
	mux.HandleFunc("GET /v1/users/search", app.searchProductHandler)

	mux.HandleFunc("POST /v1/admin/addproduct", app.addProductHandler)

	mux.HandleFunc("GET /v1/addtocart", app.addToCartHandler)
	mux.HandleFunc("GET /v1/removeitem", app.removeItemHandler)
	mux.HandleFunc("GET /v1/cartcheckout", app.cartCheckoutHandler)
	mux.HandleFunc("GET /v1/instantbuy", app.instantBuyHandler)

	return mux
}
