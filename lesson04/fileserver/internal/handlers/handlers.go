package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ptsypyshev/gb-golang-backend-level01/lesson04/fileserver/internal/fs"
)

const FilePermissions = 0600

type UploadHandler struct {
	HostAddr     string
	UploadDir    string
	TemplatesDir string
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, filepath.Join(h.TemplatesDir, "upload.html"))
	case http.MethodPost:
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}

		newFile := &fs.File{
			Link: "files/" + header.Filename,
			Name: header.Filename,
			Ext:  filepath.Ext(header.Filename),
			Size: 0,
		}

		filePath := filepath.Join(h.UploadDir, newFile.Name)

		// It is just to rename file if we have another file in 'upload' directory with the same name.
		// I think that right way is to have a model with origFileName/realFilePath relationship in DB.
		// And then use random name to save file on real FS.
		for i := 1; ; i++ {
			if fs.IsNotExist(filePath) {
				err = os.WriteFile(filePath, data, FilePermissions)
				if err != nil {
					log.Println(err)
					http.Error(w, "Unable to save file", http.StatusInternalServerError)
					return
				}
				break
			}
			newFile.Name = fs.IncrementFileName(header.Filename, i)
			newFile.Link = "files/" + newFile.Name
			filePath = filepath.Join(h.UploadDir, newFile.Name)
		}

		tmpl, err := template.ParseFiles(filepath.Join(h.TemplatesDir, "upload-success.html"))
		if err != nil {
			http.Error(w, "Unable to parse template", http.StatusBadRequest)
			return
		}
		tmplData := struct {
			Link     string
			Filename string
			Renamed  bool
		}{
			Link:     newFile.Link,
			Filename: newFile.Name,
			Renamed:  newFile.Name != header.Filename,
		}
		if err := tmpl.Execute(w, tmplData); err != nil {
			http.Error(w, fmt.Sprintf("cannot execute template: %s", err), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
}

type FilesHandler struct {
	UploadHandler
}

func (h FilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if files, err := fs.ListDir(h.UploadDir); err == nil {
			filter := r.URL.Query().Get("ext")
			if filter != "" {
				files = fs.FilterByExt(files, filter)
			}
			tmpl, err := template.ParseFiles(filepath.Join(h.TemplatesDir, "fileserver.html"))
			if err != nil {
				http.Error(w, "Unable to parse template", http.StatusBadRequest)
				return
			}

			if err := tmpl.Execute(w, files); err != nil {
				http.Error(w, fmt.Sprintf("cannot execute template: %s", err), http.StatusInternalServerError)
				return
			}
		}
	default:
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
}
