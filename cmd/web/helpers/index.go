package helpers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"runtime/debug"
	"snippetbox/cmd/web/config"

	_ "github.com/lib/pq"
)

func ServerError(app *config.Application, err error) http.HandlerFunc {
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

func ClientError(app *config.Application, status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(status), status)
	}
}

func SetupDBConn(logger *slog.Logger, dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	err = db.Ping()
	if err != nil {
		logger.Error(err.Error())
		db.Close()
		return nil
	}

	return db
}
