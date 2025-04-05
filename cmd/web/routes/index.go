package routes

import (
	"net/http"
	"snippetbox/cmd/web/handlers"
)

type Router struct {
	*http.ServeMux
}

// Initialize the routes with the application configuration.
func NewRouter(app *handlers.Application) http.Handler {

	router := &Router{
		ServeMux: http.NewServeMux(),
	}

	// Load the static files.
	app.LoadStaticFiles()(router.ServeMux)

	// Load Snippet Routes.
	router.InitSnippetRoutes(app)

	return app.RecoverPanic(app.LogRequest(app.Middlewares.CommonHeaders(router.ServeMux)))
}
