package routes

import (
	"net/http"
	"os"
	"snippetbox/cmd/web/config"
	"snippetbox/cmd/web/handlers"
)

func InitRouteHandlers(app *config.Application, addr *string) {
	router := http.NewServeMux()

	router.HandleFunc("GET /", handlers.GetSnippetHome(app))
	router.HandleFunc("GET /snippet/create", handlers.GetCreateSnippet(app))
	router.HandleFunc("POST /snippet/create", handlers.PostCreateSnippet(app))
	router.HandleFunc("PUT /snippet/update", handlers.UpdateSnippetById(app))
	router.HandleFunc("GET /snippet/view/{id}", handlers.GetSnippetById(app))

	handlers.LoadStaticFiles(app)(router)

	// Declaring the router and starting the server.
	app.Logger.Info("Starting server on %s", *addr, nil)
	err := http.ListenAndServe(*addr, router)

	app.Logger.Error(err.Error())
	os.Exit(1)
}
