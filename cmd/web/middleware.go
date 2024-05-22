package main

import (
	"fmt"
	"net/http"
	"snippetbox.pethron.me/cmd/config"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't
		// need to do this in your own code.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; fontsrc fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		// Any code here will execute on the way down the chain.
		next.ServeHTTP(w, r)
		// Any code here will execute on the way back up the chain.
	})
}

//func logRequest(app *config.Application, next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
//		next.ServeHTTP(w, r)
//	})
//}
//
//func recoverPanic(app *config.Application, next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Create a deferred function (which will always be run in the event
//		// of a panic as Go unwinds the stack).
//		defer func() {
//			// Use the builtin recover function to check if there has been a
//			// panic or not. If there has...
//			if err := recover(); err != nil {
//				// Set a "Connection: close" header on the response.
//				w.Header().Set("Connection", "close")
//				// Call the app.serverError helper method to return a 500
//				// Internal Server response.
//				app.ServerError(w, fmt.Errorf("%s", err))
//			}
//		}()
//		next.ServeHTTP(w, r)
//	})
//}

func logRequest(app *config.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
			next.ServeHTTP(w, r)
		})
	}
}

func recoverPanic(app *config.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a deferred function (which will always be run in the event
			// of a panic as Go unwinds the stack).
			defer func() {
				// Use the builtin recover function to check if there has been a
				// panic or not. If there has...
				if err := recover(); err != nil {
					// Set a "Connection: close" header on the response.
					w.Header().Set("Connection", "close")
					// Call the app.serverError helper method to return a 500
					// Internal Server response.
					app.ServerError(w, fmt.Errorf("%s", err))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
