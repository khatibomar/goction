package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/socket", http.HandlerFunc(app.socketHandler))
	mux.Post("/auct", http.HandlerFunc(app.updatePrice))

	return mux
}
