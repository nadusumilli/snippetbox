package handlers

import (
	"net/http"
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

func (app *Application) LoadStaticFiles() func(router *http.ServeMux) {
	return func(router *http.ServeMux) {
		fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
		router.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	}
}
