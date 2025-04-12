package handlers

import (
	"net/http"
	"snippetbox/cmd/web/structs"
	"snippetbox/internal/errors"
	"snippetbox/internal/models"
)

func (app *Application) UserSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := NewTemplateData[structs.UserStruct, models.User](app, r, nil, structs.UserStruct{})
		app.Render(w, r, http.StatusOK, "signup.tmpl.html", data)
	}
}

func (app *Application) UserSignupPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form *structs.UserStruct

		err := app.DecodePostForm(r, &form)
		if err != nil {
			app.ClientError(http.StatusBadRequest)(w, r)
			return
		}

		form.Validate()

		if !form.Valid() {
			data := NewTemplateData[structs.UserStruct, models.User](app, r, form, structs.UserStruct{})
			data.Form = form
			app.Render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
			return
		}

		_, err = app.Users.Insert(form.Name, form.Email, form.Password)
		if err != nil {
			if err == errors.ErrDuplicateEmail {
				form.AddFieldError("Email", "Email address is already in use")
				data := NewTemplateData[structs.UserStruct, models.User](app, r, form, structs.UserStruct{})
				app.Render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
			} else {
				app.InternalServerError(err)(w, r)
			}
			return
		}

		app.SessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

		// And redirect the user to the login page.
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func (app *Application) UserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := NewTemplateData[structs.UserLogin, structs.UserLogin](app, r, nil, structs.UserLogin{})
		app.Render(w, r, http.StatusOK, "login.tmpl.html", data)
	}
}

func (app *Application) UserLoginPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form *structs.UserLogin

		err := app.DecodePostForm(r, &form)
		if err != nil {
			app.ClientError(http.StatusBadRequest)(w, r)
			return
		}

		form.Validate()
		if !form.Valid() {
			data := NewTemplateData[structs.UserLogin, structs.UserLogin](app, r, form, structs.UserLogin{})
			data.Form = form
			app.Render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
			return
		}

		id, err := app.Users.Authenticate(form.Email, form.Password)
		if err != nil {
			if err == errors.ErrInvalidCredentials {
				form.AddNonFieldError("Email address or password is invalid")
				data := NewTemplateData[structs.UserLogin, structs.UserLogin](app, r, form, structs.UserLogin{})
				data.Form = form
				app.Render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
			} else {
				app.InternalServerError(err)(w, r)
			}
			return
		}

		app.SessionManager.Put(r.Context(), "authenticatedUserID", id)
		app.SessionManager.Put(r.Context(), "flash", "You are now logged in.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *Application) UserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the flash message from the session.
		err := app.SessionManager.RenewToken(r.Context())
		if err != nil {
			app.InternalServerError(err)(w, r)
			return
		}

		// Remove the authenticatedUserID and flash values from the session.
		app.SessionManager.Remove(r.Context(), "authenticatedUserID")
		app.SessionManager.Put(r.Context(), "flash", "You've been logged out successfully.")

		// Redirect the user to the home page.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *Application) UserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := NewTemplateData[structs.UserStruct, models.User](app, r, nil, structs.UserStruct{})
		app.Render(w, r, http.StatusOK, "profile.tmpl.html", data)
	}
}

func (app *Application) UpdateUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := NewTemplateData[structs.UserStruct, models.User](app, r, nil, structs.UserStruct{})
		app.Render(w, r, http.StatusOK, "update_profile.tmpl.html", data)
	}
}

func (app *Application) ResetUserPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := NewTemplateData[structs.UserStruct, models.User](app, r, nil, structs.UserStruct{})
		app.Render(w, r, http.StatusOK, "update_profile.tmpl.html", data)
	}
}
