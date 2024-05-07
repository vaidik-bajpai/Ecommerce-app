package main

import (
	"errors"
	"net/http"

	"github.com/vaidik-bajpai/ecommerce-api/internal/data"
	"github.com/vaidik-bajpai/ecommerce-api/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Phone     string `json:"phone"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		FirstName: &input.FirstName,
		LastName:  &input.LastName,
		Email:     &input.Email,
		Phone:     &input.Phone,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	if data.ValidateUser(v, user); len(v.Errors) != 0 {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	token := "aksdkabdja"
	user.Token = &token
	user.RefreshToken = &token

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddErrors("email", "a user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicatePhoneNo):
			v.AddErrors("phone", "a user with this phone no. already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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
