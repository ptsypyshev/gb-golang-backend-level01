package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ptsypyshev/gb-golang-backend-level01/lesson04/fileserver/internal/handlers"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		cwd, _ = os.UserHomeDir()
	}

	uploadHandler := &handlers.UploadHandler{
		HostAddr:     "127.0.0.1:80",
		UploadDir:    filepath.Join(cwd, "upload"),
		TemplatesDir: filepath.Join(cwd, "templates"),
	}

	filesHandler := &handlers.FilesHandler{
		UploadHandler: *uploadHandler,
	}

	dirToServe := http.Dir(uploadHandler.UploadDir)
	fs := &http.Server{
		Addr:         uploadHandler.HostAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.Handle("/", filesHandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(dirToServe)))
	http.Handle("/upload/", uploadHandler)

	log.Fatal(fs.ListenAndServe())
}
