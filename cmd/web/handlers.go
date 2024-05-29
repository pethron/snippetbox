package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"snippetbox.pethron.me/cmd/config"
	"snippetbox.pethron.me/internal/models"
	"snippetbox.pethron.me/internal/validator"
	"strconv"
)

func home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.ServerError(w, err)
			return
		}

		data := app.NewTemplateData(r)
		data.Snippets = snippets

		app.Render(w, http.StatusOK, "home.tmpl", data)
	}
}

func snippetView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())

		id, err := strconv.Atoi(params.ByName("id"))
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

		data := app.NewTemplateData(r)
		data.Snippet = snippet

		app.Render(w, http.StatusOK, "view.tmpl", data)
	}
}

func snippetCreate(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.NewTemplateData(r)
		data.Form = snippetCreateForm{
			Expires: 365,
		}

		app.Render(w, http.StatusOK, "create.tmpl", data)
	}
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func snippetCreatePost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form snippetCreateForm

		err := app.DecodePostForm(r, &form)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

		if !form.Valid() {
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}

		id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			app.ServerError(w, err)
		}

		app.SessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func userSignup(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.NewTemplateData(r)
		data.Form = userSignupForm{}
		app.Render(w, http.StatusOK, "signup.tmpl", data)
	}
}

func userSignupPost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form userSignupForm

		err := app.DecodePostForm(r, &form)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
		form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be more than 100 characters long")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
		form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

		if !form.Valid() {
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		}

		err = app.Users.Insert(form.Name, form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				form.AddFieldError("email", "Email address is already in use")
				data := app.NewTemplateData(r)
				data.Form = form
				app.Render(w, http.StatusUnprocessableEntity, "signup.tmpl",
					data)
			} else {
				app.ServerError(w, err)
			}
			return
		}
		// Otherwise add a confirmation flash message to the session confirming that
		// their signup worked.
		app.SessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
		// And redirect the user to the login page.
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	}
}

// Create a new userLoginForm struct.
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func userLogin(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.NewTemplateData(r)
		data.Form = userLoginForm{}
		app.Render(w, http.StatusOK, "login.tmpl", data)
	}
}

func userLoginPost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form userLoginForm

		err := app.DecodePostForm(r, &form)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

		if !form.Valid() {
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
			return
		}

		id, err := app.Users.Authenticate(form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.AddNonFieldError("Email or password is incorrect")
				data := app.NewTemplateData(r)
				data.Form = form
				app.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
			} else {
				app.ServerError(w, err)
			}
		}

		err = app.SessionManager.RenewToken(r.Context())
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.SessionManager.Put(r.Context(), "authenticatedUserID", id)

		redirect := app.SessionManager.PopString(r.Context(), "redirectPathAfterLogin")
		if redirect != "" {
			r.URL.Path = redirect
			http.Redirect(w, r, redirect, http.StatusSeeOther)
		}

		http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	}
}

func userLogoutPost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.SessionManager.RenewToken(r.Context())
		if err != nil {
			app.ServerError(w, err)
			return
		}
		app.SessionManager.Remove(r.Context(), "authenticatedUserID")
		app.SessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func about(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.NewTemplateData(r)
		app.Render(w, http.StatusOK, "about.tmpl", data)
	}
}

func accountView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := app.SessionManager.GetInt(r.Context(), "authenticatedUserID")

		user, err := app.Users.Get(userID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFoundError(w)
			} else {
				app.ServerError(w, err)
			}
			return
		}

		data := app.NewTemplateData(r)
		data.User = user
		app.Render(w, http.StatusOK, "account.tmpl", data)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
