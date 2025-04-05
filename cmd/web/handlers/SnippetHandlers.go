package handlers

import (
	"errors"
	"fmt"
	"net/http"
	structs "snippetbox/cmd/web/structs/snippets"
	"snippetbox/cmd/web/templates"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"
	"time"
)

func (app *Application) GetCreateSnippet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.TemplateData{
			CurrentYear: time.Now().Year(),
		}

		app.Render(w, r, http.StatusOK, "create.tmpl.html", &data)
	}
}

func (app *Application) PostCreateSnippet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var form structs.SnippetStruct

		err := app.DecodePostForm(r, &form)
		if err != nil {
			app.ClientError(http.StatusBadRequest)(w, r)
			return
		}

		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

		if !form.Valid() {
			data := &templates.TemplateData{
				CurrentYear: time.Now().Year(),
				Form:        form,
			}
			app.Render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
			return
		}

		id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			app.InternalServerError(err)(w, r)
			return
		}

		// Use the Put() method to add a string value ("Snippet successfully
		// created!") and the corresponding key ("flash") to the session data.
		app.SessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}

func (app *Application) GetSnippetHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.InternalServerError(err)(w, r)
			return
		}

		app.Render(w, r, http.StatusOK, "home.tmpl.html", &templates.TemplateData{
			CurrentYear: time.Now().Year(),
			Snippets:    snippets,
		})
	}
}

func (app *Application) UpdateSnippetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			app.NotFound(err)(w, r)
			return
		}

		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.InternalServerError(err)(w, r)
			return
		}

		app.Render(w, r, http.StatusOK, "home.tmpl.html", &templates.TemplateData{
			Snippets: snippets,
		})
	}
}

func (app *Application) GetSnippetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil || id < 1 {
			app.NotFound(err)(w, r)
			return
		}

		snippet, err := app.Snippets.Get(id)

		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(err)(w, r)
			} else {
				app.InternalServerError(err)(w, r)
			}
			return
		}

		flash := app.SessionManager.PopString(r.Context(), "flash")

		data := &templates.TemplateData{
			CurrentYear: time.Now().Year(),
			Snippet:     snippet,
			Flash:       flash,
		}

		app.Render(w, r, http.StatusOK, "view.tmpl.html", data)
	}
}

func (app *Application) GetAllSnippets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.InternalServerError(err)(w, r)
			return
		}

		fmt.Fprintf(w, "%+v", snippets)
	}
}
