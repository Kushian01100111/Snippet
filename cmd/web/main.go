package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	server := http.NewServeMux()

	// -> Problemas con el file server, debe de devolver 404 error en vez de 500 error
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})

	server.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	server.HandleFunc("GET /{$}", Home)
	server.HandleFunc("GET /snippet/view/{id}", snippetView)
	server.HandleFunc("GET /snippet/create", snippetCreate)
	server.HandleFunc("POST /snippet/create", snippetCreatePost)

	logger.Info("Startring server", "addr", *addr)

	err := http.ListenAndServe(*addr, server)
	logger.Error(err.Error())
	os.Exit(1)
}

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
