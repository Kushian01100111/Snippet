package main

import (
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (app *application) routes() *http.ServeMux {
	server := http.NewServeMux()

	// -> Problemas con el file server, debe de devolver 404 error en vez de 500 error
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	server.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	server.HandleFunc("GET /{$}", app.Home)
	server.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	server.HandleFunc("GET /snippet/create", app.snippetCreate)
	server.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return server
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
