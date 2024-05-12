package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/users/signup", app.createUserHandler)
	mux.HandleFunc("POST /v1/users/login", app.userLoginHandler)
	mux.HandleFunc("POST /v1/users/productview", app.authenticate(app.viewProductHandler))

	mux.HandleFunc("GET /v1/products/search", app.searchProductHandler)
	mux.HandleFunc("GET /v1/products/{id}", app.searchProductByIDHandler)

	mux.HandleFunc("POST /v1/admin/products", app.createProductHandler)
	mux.HandleFunc("DELETE /v1/admin/products/{id}", app.removeItemHandler)

	mux.HandleFunc("GET /v1/addtocart", app.addToCartHandler)
	mux.HandleFunc("GET /v1/removeitem", app.removeItemHandler)
	mux.HandleFunc("GET /v1/cartcheckout", app.cartCheckoutHandler)
	mux.HandleFunc("GET /v1/instantbuy", app.instantBuyHandler)

	return mux
}
