package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/pascaldekloe/jwt"
	"github.com/vaidik-bajpai/ecommerce-api/internal/data"
	"github.com/vaidik-bajpai/ecommerce-api/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	app.models.Users.Get(6)
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
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			app.invalidCredentialResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if !match {
		app.invalidCredentialResponse(w, r)
		return
	}

	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(user.ID, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "bajpai.ecommerce"
	claims.Audiences = []string{"bajpai.ecommerce"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": string(jwtBytes)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) viewProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("viewProduct endpoint"))
}
