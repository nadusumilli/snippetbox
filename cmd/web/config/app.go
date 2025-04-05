package config

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"runtime/debug"
	"snippetbox/cmd/web/constants"
	"snippetbox/cmd/web/middlewares"
	"snippetbox/cmd/web/templates"
	"snippetbox/internal/models"
)

type ApplicationConfig struct {
	Logger        *slog.Logger
	Middlewares   *middlewares.Middlewares
	DB            *sql.DB
	Snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func NewApplicationConfigConnection(logger *slog.Logger, dsn *string) *ApplicationConfig {
	// Initializing the database connection and module.
	db := NewDatabaseConnection(dsn, logger)

	// initialize the template cache
	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		logger.Error("Error creating template cache", "error", err)
		return nil
	}

	// Initialize the config with the database connection, logger and snippets model and return the instance.
	return &ApplicationConfig{
		Logger:        logger,
		DB:            db.DB,
		Middlewares:   middlewares.NewMiddlewares(),
		Snippets:      models.NewSnippetModel(db.DB),
		templateCache: templateCache,
	}
}

func (app *ApplicationConfig) InternalServerError(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			method = r.Method
			uri    = r.URL.RequestURI()
			trace  = string(debug.Stack())
		)

		app.Logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *ApplicationConfig) NotFound(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			method = r.Method
			uri    = r.URL.RequestURI()
			trace  = string(debug.Stack())
		)
		app.Logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (app *ApplicationConfig) BadRequest(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			method = r.Method
			uri    = r.URL.RequestURI()
			trace  = string(debug.Stack())
		)
		app.Logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (app *ApplicationConfig) ClientError(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(status), status)
	}
}

func (app *ApplicationConfig) Render(w http.ResponseWriter, r *http.Request, status int, templateName string, data *templates.TemplateData) {
	// get the template from the cache
	tmpl, ok := app.templateCache[templateName]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", templateName)
		app.InternalServerError(err)(w, r)
		return
	}

	buf := new(bytes.Buffer)

	// Execute the template set and write the response body. Again, if there
	// is any error we call the serverError() helper.
	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.InternalServerError(err)(w, r)
		return
	}

	// Write out the provided HTTP status code ('200 OK', '400 Bad Request' etc).
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *ApplicationConfig) SuccessResponseWriter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Create an instance of the SuccessResponse struct
		successResponse := constants.SuccessResponse{
			Message: "Success",
			SDESC:   "Request was successful",
			SCODE:   "200",
		}
		// Marshal the struct to JSON
		jsonResponse, err := json.Marshal(successResponse)
		if err != nil {
			app.Logger.Error("Error marshalling JSON", "error", err)
			app.InternalServerError(err)(w, r)
			return
		}
		// Write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

func (app *ApplicationConfig) ErrorResponseWriter(data map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

		result := make(map[string]interface{})
		for key, value := range data {
			result[key] = value
		}

		// Create an instance of the ErrorResponse struct
		errorResponse := constants.ErrorResponse{
			Error: "Internal Server Error",
			SDESC: "An unexpected error occurred",
			SCODE: "500",
			DATA:  result,
		}

		// Marshal the struct to JSON
		jsonResponse, err := json.Marshal(errorResponse)
		if err != nil {
			app.Logger.Error("Error marshalling JSON", "error", err)
			app.InternalServerError(err)(w, r)
			return
		}

		// Write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
	}
}

func (app *ApplicationConfig) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.Logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (app *ApplicationConfig) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic error
				app.Logger.Error("panic error", "error", err, "stack", string(debug.Stack()))

				// Send a 500 Internal Server Error response
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
