package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"auth-mail/internal/data"
	"auth-mail/internal/validator"
)

// registerUserInput defines expected JSON payload
type registerUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// registerUserHandler handles POST /v1/users
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input registerUserInput

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Username:  input.Username,
		Email:     input.Email,
		Activated: false,
	}

	if err := user.Password.Set(input.Password); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	data.ValidateUser(v, user)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Users.Insert(user); err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			v.AddError("email", "a user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Generate activation token
	token, err := data.GenerateToken(user.ID, 3*24*time.Hour, "activation")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Tokens.Insert(token)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send email (Mailtrap)
	app.background(func() {
		err := app.mailer.Send(
			user.Email,
			"Activate your account",
			"Your activation token: "+token.Plaintext,
		)
		if err != nil {
			log.Println(err)
		}
	})

	// Response
	err = app.writeJSON(w, http.StatusCreated, envelope{
		"user": user,
	}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}