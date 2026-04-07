package main

import (
	"errors"
	"net/http"
	"time"

	"auth-mail/internal/data"
	"auth-mail/internal/validator"
)

type createAuthenticationTokenInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input createAuthenticationTokenInput

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.invalidCredentialsResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match || !user.Activated {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := data.GenerateToken(user.ID, 24*time.Hour, "authentication")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Tokens.Insert(token)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{
		"authentication_token": token,
	}, nil)
}