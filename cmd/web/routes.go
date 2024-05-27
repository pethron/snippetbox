package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"snippetbox.pethron.me/cmd/config"
)

func routes(app *config.Application) func() http.Handler {
	return func() http.Handler {
		router := httprouter.New()

		router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.NotFoundError(w)
		})

		fileServer := http.FileServer(http.Dir("./ui/static/"))

		dynamic := alice.New(app.SessionManager.LoadAndSave)

		router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

		router.Handler(http.MethodGet, "/", dynamic.ThenFunc(home(app)))
		router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(snippetView(app)))
		router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(snippetCreate(app)))
		router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(snippetCreatePost(app)))

		standard := alice.New(recoverPanic(app), logRequest(app), secureHeaders)
		return standard.Then(router)

		//return recoverPanic(app)(logRequest(app)(secureHeaders(mux)))
	}
}
