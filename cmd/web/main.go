package main

import (
	"crypto/tls"
	"flag"
	"log/slog"
	"net/http"
	"snippetbox/cmd/web/constants"
	"snippetbox/cmd/web/handlers"
	"snippetbox/cmd/web/routes"
	"time"
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
	router := routes.InitRoutes(app)

	// TLS configuration.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Start the server.
	server := &http.Server{
		Addr:         *addr,
		Handler:      router,
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	listenErr := server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	if listenErr != nil {
		app.Logger.Error("Failed to start the server", "error", listenErr)
	}

	// Log the server start.
	app.Logger.Info("Server started", "address", *addr)
}
