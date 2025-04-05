package config

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"snippetbox/cmd/web/constants"
)

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
