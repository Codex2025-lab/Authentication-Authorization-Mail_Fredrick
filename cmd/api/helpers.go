package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type envelope map[string]any

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) background(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic recovered:", err)
			}
		}()
		fn()
	}()
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "server error"}, nil)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.writeJSON(w, http.StatusUnprocessableEntity, envelope{"errors": errors}, nil)
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusUnauthorized, envelope{
		"error": "invalid authentication credentials",
	}, nil)
}