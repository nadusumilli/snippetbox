package handlers

import (
	"log/slog"
	"net/http"
	"os"
	"snippetbox/cmd/web/config"
	"snippetbox/cmd/web/templates"
	"snippetbox/internal/validator"
	"time"
)

type Application struct {
	*config.ApplicationConfig
}

func NewApiConnection(dsn *string, addr *string) *Application {

	// Initializing the structured logger module.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// Initializing the application configuration module.
	appConfig := config.NewApplicationConfigConnection(logger, dsn)
	if appConfig == nil {
		logger.Error("Error creating application config connection")
		os.Exit(1)
	}

	app := &Application{
		ApplicationConfig: appConfig,
	}

	// returning the app.
	return app
}

func NewTemplateData[T any, M any](app *Application, r *http.Request, form *T, model T) *templates.TemplateData[T, M] {
	if form == nil {
		form = new(T)
		if v, ok := any(form).(interface{ SetValidator(v validator.Validator) }); ok {
			v.SetValidator(validator.New(model))
		}
	}

	return &templates.TemplateData[T, M]{
		CurrentYear: time.Now().Year(),
		Flash:       app.SessionManager.PopString(r.Context(), "flash"),
		Form:        form,
	}
}
