package config

import (
	"github.com/justinas/nosurf"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"snippetbox.pethron.me/internal/models"
	"snippetbox.pethron.me/ui"
	"time"
)

type TemplateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {

		return nil, err
	}

	// Loop through the page filepaths one-by-one.
	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files,
			patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}

func (app *Application) NewTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.SessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}
