package routes

import (
	"net/http"
	"snippetbox/cmd/web/handlers"
)

func NewSnippetRouter(app *handlers.Application) http.Handler {
	r := NewRouter()
	InitSnippetRoutes(r, app)
	return app.SessionManager.LoadAndSave(r.Handler())
}

func InitSnippetRoutes(r *Router, app *handlers.Application) {
	r.HandleFunc("GET /latest", app.GetSnippetHome())
	r.HandleFunc("GET /create", app.GetCreateSnippet())
	r.HandleFunc("POST /create", app.PostCreateSnippet())
	r.HandleFunc("PUT /update", app.UpdateSnippetById())
	r.HandleFunc("GET /view/{id}", app.GetSnippetById())
	r.HandleFunc("GET /list", app.GetAllSnippets())
}
