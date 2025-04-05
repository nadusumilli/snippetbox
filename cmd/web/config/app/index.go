package config

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
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
	formDecoder   *form.Decoder
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
		formDecoder:   form.NewDecoder(),
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

func (app *ApplicationConfig) DecodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err := app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}
