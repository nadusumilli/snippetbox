package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"snippetbox/cmd/web/config"
	"snippetbox/cmd/web/helpers"
	"snippetbox/cmd/web/templates"
	"snippetbox/internal/models"
	"strconv"
)

func GetCreateSnippet(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Create a new snippet.")
	}
}

func PostCreateSnippet(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new snippet post.
		title := "snail"
		content := "A snail is a small animal with a soft body that moves very slowly and has a spiral-shaped shell on its back."
		expires := 7

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			helpers.ServerError(app, err)(w, r)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}

func GetSnippetHome(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		files := []string{
			"./ui/html/partials/nav.tmpl.html",
			"./ui/html/pages/base.tmpl.html",
			"./ui/html/pages/home.tmpl.html",
		}

		ts, err := template.ParseFiles(files...) // The path should either be relative to the root of the project or absolute path to file.
		if err != nil {
			helpers.ServerError(app, err)(w, r)
			return
		}

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			helpers.ServerError(app, err)(w, r)
		}
	}
}

func UpdateSnippetById(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Updating the Snippet with id: %d", id)
	}
}

func GetSnippetById(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}

		snippet, err := app.Snippets.Get(id)

		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				helpers.ServerError(app, err)(w, r)
			}
			return
		}

		fmt.Fprintf(w, "%+v", snippet)
	}
}

func GetAllSnippets(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.Snippets.Latest()
		if err != nil {
			helpers.ServerError(app, err)(w, r)
			return
		}

		fmt.Fprintf(w, "%+v", snippets)
	}
}

func SnippetView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				helpers.ServerError(app, err)
			}
			return
		}

		files := []string{
			"./ui/html/partials/nav.tmpl.html",
			"./ui/html/pages/base.tmpl.html",
			"./ui/html/pages/view.tmpl.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			helpers.ServerError(app, err)
			return
		}

		// Create an instance of a templateData struct holding the snippet data.
		data := templates.TemplateData{
			Snippet: snippet,
		}

		// Pass in the templateData struct when executing the template.
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			helpers.ServerError(app, err)
		}
	}
}
