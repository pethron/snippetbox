package main

import (
	"github.com/justinas/alice"
	"net/http"
	"snippetbox.pethron.me/cmd/config"
)

func routes(app *config.Application) func() http.Handler {
	return func() http.Handler {
		mux := http.NewServeMux()

		fileServer := http.FileServer(http.Dir("./ui/static/"))
		mux.Handle("/static/", http.StripPrefix("/static", fileServer))

		mux.HandleFunc("/", home(app))
		mux.HandleFunc("/snippet/view", snippetView(app))
		mux.HandleFunc("/snippet/create", snippetCreate(app))

		standard := alice.New(recoverPanic(app), logRequest(app), secureHeaders)
		return standard.Then(mux)

		//return recoverPanic(app)(logRequest(app)(secureHeaders(mux)))
	}
}
