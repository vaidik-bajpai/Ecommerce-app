package main

import "net/http"

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("createUser endpoint"))
}

func (app *application) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("userLogin endpoint"))
}

func (app *application) viewProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("viewProduct endpoint"))
}

func (app *application) searchProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("searchProduct endpoint"))
}
