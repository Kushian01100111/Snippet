package main

import (
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (app *application) routes() http.Handler {
	server := http.NewServeMux()

	// -> Problemas con el file server, debe de devolver 404 error en vez de 500 error
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	server.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	server.Handle("GET /{$}", dynamic.ThenFunc(app.Home))
	server.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	server.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	server.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	server.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	server.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	server.Handle("GET /user/login", dynamic.ThenFunc(app.userLoginin))
	server.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	server.Handle("POST /user/logout", dynamic.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(server)
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
