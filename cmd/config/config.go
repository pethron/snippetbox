package config

import (
	"log"
	"snippetbox.pethron.me/internal/models"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Snippets *models.SnippetModel
}
