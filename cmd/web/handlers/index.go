package handlers

import (
	"log/slog"
	"os"
	"snippetbox/cmd/web/config"
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
