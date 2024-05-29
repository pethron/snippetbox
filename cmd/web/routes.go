package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"snippetbox.pethron.me/cmd/config"
	"snippetbox.pethron.me/ui"
)

func routes(app *config.Application) func() http.Handler {
	return func() http.Handler {
		router := httprouter.New()

		router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.NotFoundError(w)
		})

		fileServer := http.FileServer(http.FS(ui.Files))

		router.Handler(http.MethodGet, "/static/*filepath", fileServer)

		router.HandlerFunc(http.MethodGet, "/ping", ping)

		// unprotected
		dynamic := alice.New(app.SessionManager.LoadAndSave, noSurf, Authenticate(app))

		router.Handler(http.MethodGet, "/", dynamic.ThenFunc(home(app)))
		router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(snippetView(app)))
		router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(userSignup(app)))
		router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(userSignupPost(app)))
		router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(userLogin(app)))
		router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(userLoginPost(app)))
		router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(about(app)))

		// protected

		protected := dynamic.Append(requireAuthentication(app))

		router.Handler(http.MethodGet, "/account/view", protected.ThenFunc(accountView(app)))
		router.Handler(http.MethodGet, "/account/password/update", protected.ThenFunc(accountPasswordUpdate(app)))
		router.Handler(http.MethodPost, "/account/password/update", protected.ThenFunc(accountPasswordUpdatePost(app)))
		router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(snippetCreate(app)))
		router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(snippetCreatePost(app)))
		router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(userLogoutPost(app)))

		standard := alice.New(recoverPanic(app), logRequest(app), secureHeaders)
		return standard.Then(router)

		//return recoverPanic(app)(logRequest(app)(secureHeaders(mux)))
	}
}
