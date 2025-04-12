package routes

import (
	"net/http"
	"snippetbox/cmd/web/handlers"
)

func NewUserRouter(app *handlers.Application) http.Handler {
	r := NewRouter()
	InitUserRoutes(r, app)
	return app.SessionManager.LoadAndSave(r.Handler())
}

func InitUserRoutes(r *Router, app *handlers.Application) {
	r.HandleFunc("GET /signup", app.UserSignup())
	r.HandleFunc("POST /signup", app.UserSignupPost())
	r.HandleFunc("GET /login", app.UserLogin())
	r.HandleFunc("POST /login", app.UserLoginPost())
	r.HandleFunc("POST /logout", app.UserLogout())
	r.HandleFunc("GET /profile", app.UserProfile())
	r.HandleFunc("PUT /profile/update", app.UpdateUserProfile())
	r.HandleFunc("POST /reset-password", app.ResetUserPassword())
}
