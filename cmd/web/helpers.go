package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"snippetbox.pethron.me/cmd/config"
)

func serverError(app *config.Application) func(w http.ResponseWriter, err error) {
	return func(w http.ResponseWriter, err error) {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func clientError(app *config.Application) func(w http.ResponseWriter, status int) {
	return func(w http.ResponseWriter, status int) {
		http.Error(w, http.StatusText(status), status)
	}
}

func notFoundError(app *config.Application) func(w http.ResponseWriter) {
	return func(w http.ResponseWriter) {
		clientError(app)(w, http.StatusNotFound)
	}
}
