package main

import (
	"flag"
	"net/http"
	"snippetbox/cmd/web/constants"
	"snippetbox/cmd/web/handlers"
	"snippetbox/cmd/web/routes"
)

func main() {

	// Getting the address from the command line flag.
	addr := flag.String("addr", constants.PORT, "HTTP network address")
	dsn := flag.String("dsn", constants.DATABASE_CONNECTION_STRING, "PostgreSQL data source name")

	// Parsing the command line flags.
	flag.Parse()

	// Initialize the app config, logger and database.
	app := handlers.NewApiConnection(dsn, addr)

	// Derer the closing of the database if application closes.
	defer func() {
		if err := app.DB.Close(); err != nil {
			app.Logger.Error("Failed to close the database connection", "error", err)
		}
	}()

	// initialize the routes for our api's.
	router := routes.NewRouter(app)

	// Start the server.
	listenErr := http.ListenAndServe(*addr, router)

	if listenErr != nil {
		app.Logger.Error("Failed to start the server", "error", listenErr)
	}

	// Log the server start.
	app.Logger.Info("Server started", "address", *addr)
}
