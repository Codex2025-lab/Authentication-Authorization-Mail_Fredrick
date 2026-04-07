package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/v1/users", app.registerUserHandler)
	router.HandleFunc("/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return router
}