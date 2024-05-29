package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
	"snippetbox.pethron.me/internal/models"
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Users          models.UserModelInterface
	Snippets       models.SnippetModelInterface
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
	Debug          bool
}
