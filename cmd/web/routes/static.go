package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"snippetbox/cmd/web/handlers"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}

func NewStaticRouter(app *handlers.Application) http.Handler {
	r := NewRouter()
	InitStaticRoutes(r, app)
	return app.SessionManager.LoadAndSave(r.Handler())
}

func InitStaticRoutes(r *Router, app *handlers.Application) {
	cwd, err := os.Getwd()
	if err != nil {
		app.Logger.Error("Failed to get current working directory", err)
		return
	}

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(filepath.Join(cwd, "ui", "static"))})
	app.Logger.Info("Loading static files from", "file path:", fileServer)
	r.Handle("GET /", fileServer)
}
