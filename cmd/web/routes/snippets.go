package routes

import "snippetbox/cmd/web/handlers"

func (router *Router) InitSnippetRoutes(app *handlers.Application) {
	router.HandleFunc("GET /", app.GetSnippetHome())
	router.HandleFunc("GET /snippet/create", app.GetCreateSnippet())
	router.HandleFunc("POST /snippet/create", app.PostCreateSnippet())
	router.HandleFunc("PUT /snippet/update", app.UpdateSnippetById())
	router.HandleFunc("GET /snippet/view/{id}", app.SnippetView())
	router.HandleFunc("GET /snippets", app.GetAllSnippets())
}
