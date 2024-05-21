package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox.pethron.me/cmd/config"
	"snippetbox.pethron.me/internal/models"
	"strconv"
)

func home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.NotFoundError(w)
			return
		}

		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.ServerError(w, err)
			return
		}

		for _, snippet := range snippets {
			fmt.Fprintf(w, "%+v\n", snippet)
		}

		//files := []string{
		//	"./ui/html/base.tmpl",
		//	"./ui/html/pages/home.tmpl",
		//	"./ui/html/partials/nav.tmpl",
		//}
		//
		//ts, err := template.ParseFiles(files...)
		//if err != nil {
		//	app.ErrorLog.Print(err.Error())
		//	app.ServerError(w, err)
		//	return
		//}
		//
		//err = ts.ExecuteTemplate(w, "base", nil)
		//if err != nil {
		//	app.ErrorLog.Print(err.Error())
		//	app.ServerError(w, err)
		//}
	}
}

func snippetView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.NotFoundError(w)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFoundError(w)
			} else {
				app.ServerError(w, err)
			}
			return
		}

		fmt.Fprintf(w, "%v", snippet)
	}
}

func snippetCreate(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.ClientError(w, http.StatusMethodNotAllowed)
			return
		}

		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n–Kobayashi Issa"
		expires := 7
		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			app.ServerError(w, err)
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	}
}
