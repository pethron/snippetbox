package main

import (
	"net/http"
	"snippetbox.pethron.me/cmd/config"
)

func routes(app *config.Application) func() *http.ServeMux {
	return func() *http.ServeMux {
		mux := http.NewServeMux()

		fileServer := http.FileServer(http.Dir("./ui/static/"))
		mux.Handle("/static/", http.StripPrefix("/static", fileServer))

		mux.HandleFunc("/", home(app))
		mux.HandleFunc("/snippet/view", snippetView(app))
		mux.HandleFunc("/snippet/create", snippetCreate(app))

		return mux
	}
}
