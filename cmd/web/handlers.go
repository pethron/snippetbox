package main

import (
	"fmt"
	"html/template"
	"net/http"
	"snippetbox.pethron.me/cmd/config"
	"strconv"
)

func home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFoundError(app)
			return
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/pages/home.tmpl",
			"./ui/html/partials/nav.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Print(err.Error())
			serverError(app)(w, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			app.ErrorLog.Print(err.Error())
			serverError(app)(w, err)
		}
	}
}

func snippetView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			notFoundError(app)
			return
		}
		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}

func snippetCreate(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			clientError(app)(w, http.StatusMethodNotAllowed)
			return
		}

		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n–Kobayashi Issa"
		expires := 7
		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			serverError(app)(w, err)
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	}
}
