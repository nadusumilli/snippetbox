package main

import (
	"flag"
	"log/slog"
	"os"
	"snippetbox/cmd/web/config"
	"snippetbox/cmd/web/helpers"
	"snippetbox/cmd/web/routes"
	"snippetbox/internal/models"
)

func main() {
	// Database connectiong string for local.
	connString := "user=web password=snippet@123 dbname=snippetbox sslmode=disable"

	// Getting the address from the command line flag.
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Getting the database connection from the command line flag.
	dsn := flag.String("dsn", connString, "PostgreSQL data source name")
	flag.Parse()

	// Creating a new logger.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// Setting up the sql connections.
	db := helpers.SetupDBConn(logger, *dsn)
	defer db.Close()

	// Ensuring that the dependencies are passed properly for all the handlers.
	// Adding loggers and models to the application struct.
	app := &config.Application{
		Logger:   logger,
		Snippets: &models.SnippetModel{DB: db},
	}

	routes.InitRouteHandlers(app, addr)
}
