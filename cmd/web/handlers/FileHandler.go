package handlers

import (
	"net/http"
	"os"
	"path/filepath"
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

func (app *Application) LoadStaticFiles(router *http.ServeMux) {
	cwd, err := os.Getwd()
	if err != nil {
		app.Logger.Error("Failed to get current working directory", err)
		return
	}

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(filepath.Join(cwd, "../../ui", "static"))})
	app.Logger.Info("Loading static files from", "file path:", fileServer)
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))
}
