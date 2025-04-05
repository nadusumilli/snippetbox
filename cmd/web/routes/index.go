package routes

import (
	"net/http"
	"snippetbox/cmd/web/handlers"
)

type Router struct {
	*http.ServeMux
}

func NewRouter(app *handlers.Application) http.Handler {
	sessionRouter := &Router{ServeMux: http.NewServeMux()}
	sessionRouter.InitSnippetRoutes(app)

	staticRouter := &Router{ServeMux: http.NewServeMux()}
	app.LoadStaticFiles(staticRouter.ServeMux)

	masterMux := http.NewServeMux()
	masterMux.Handle("/static/", staticRouter.ServeMux)
	masterMux.Handle("/", app.SessionManager.LoadAndSave(sessionRouter.ServeMux))

	return app.RecoverPanic(
		app.LogRequest(
			app.Middlewares.CommonHeaders(masterMux),
		),
	)
}
