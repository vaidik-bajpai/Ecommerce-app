package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/users", app.createUserHandler)
	mux.HandleFunc("POST /v1/users/login", app.userLoginHandler)

	mux.HandleFunc("GET /v1/products", app.searchProductHandler)
	mux.HandleFunc("GET /v1/products/{id}", app.searchProductByIDHandler)

	mux.HandleFunc("POST /v1/admin/products", app.requireAuthenticatedUser(app.createProductHandler))
	mux.HandleFunc("DELETE /v1/admin/products/{id}", app.requireAuthenticatedUser(app.deleteProductHandler))

	mux.HandleFunc("GET /v1/addtocart", app.requireAuthenticatedUser(app.addToCartHandler))
	mux.HandleFunc("GET /v1/removefromcart", app.requireAuthenticatedUser(app.removeItemHandler))

	mux.HandleFunc("GET /v1/cartcheckout", app.requireAuthenticatedUser(app.cartCheckoutHandler))
	mux.HandleFunc("GET /v1/instantbuy", app.requireAuthenticatedUser(app.instantBuyHandler))

	return app.authenticate(mux)
}
