package main

import "net/http"

func (app *application) addProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("addProduct endpoint"))
}
