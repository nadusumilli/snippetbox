package config

import (
	"log/slog"
	"snippetbox/internal/models"
)

type Application struct {
	Logger   *slog.Logger
	Snippets *models.SnippetModel
}
