package routes

import (
	"net/http"
	"snippetbox/cmd/web/handlers"
)

type Router struct {
	Mux *http.ServeMux
}

func InitRoutes(app *handlers.Application) http.Handler {
	masterMux := http.NewServeMux()
	masterMux.Handle("/static/", http.StripPrefix("/static", NewStaticRouter(app)))
	masterMux.Handle("/user/", http.StripPrefix("/user", NewUserRouter(app)))
	masterMux.Handle("/snippet/", http.StripPrefix("/snippet", NewSnippetRouter(app)))
	masterMux.Handle("/",
		app.SessionManager.LoadAndSave(app.GetSnippetHome()))

	return app.RecoverPanic(
		app.LogRequest(
			app.Middlewares.CommonHeaders(masterMux),
		),
	)
}

func NewRouter() *Router {
	return &Router{Mux: http.NewServeMux()}
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	r.Mux.Handle(pattern, handler)
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.Mux.HandleFunc(pattern, handler)
}

func (r *Router) Handler() http.Handler {
	return r.Mux
}
