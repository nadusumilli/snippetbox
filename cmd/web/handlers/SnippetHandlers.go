package handlers

import (
	"errors"
	"fmt"
	"net/http"
	structs "snippetbox/cmd/web/structs"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"
)

func (app *Application) GetCreateSnippet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := NewTemplateData[structs.SnippetStruct, models.Snippet](app, r, nil, structs.SnippetStruct{})

		app.Render(w, r, http.StatusOK, "create.tmpl.html", data)
	}
}

func (app *Application) PostCreateSnippet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var form = &structs.SnippetStruct{
			Validator: validator.New(&models.Snippet{}),
		}

		err := app.DecodePostForm(r, &form)
		if err != nil {
			app.ClientError(http.StatusBadRequest)(w, r)
			return
		}

		form.Validate()

		if !form.Valid() {
			data := NewTemplateData[structs.SnippetStruct, models.Snippet](app, r, nil, structs.SnippetStruct{})
			data.Form = form
			app.Render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
			return
		}

		id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			app.InternalServerError(err)(w, r)
			return
		}

		// Set a flash message to the session
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

		data := NewTemplateData[structs.SnippetStruct, []models.Snippet](app, r, nil, structs.SnippetStruct{})
		data.Data = snippets

		app.Render(w, r, http.StatusOK, "home.tmpl.html", data)
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

		data := NewTemplateData[structs.SnippetStruct, []models.Snippet](app, r, nil, structs.SnippetStruct{})
		data.Data = snippets

		app.Render(w, r, http.StatusOK, "home.tmpl.html", data)
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

		data := NewTemplateData[structs.SnippetStruct, models.Snippet](app, r, nil, structs.SnippetStruct{})
		data.Data = snippet

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
