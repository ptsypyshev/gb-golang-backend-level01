package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func createTestFile(fullFilePath string) error {
	content := []byte("test\n")
	if err := os.WriteFile(fullFilePath, content, 0644); err != nil {
		return err
	}
	return nil
}

func TestUploadHandler_ServeHTTP(t *testing.T) {
	cwd := filepath.Join(getCurrentPath(), "..", "..")

	testFilePath := filepath.Join(os.TempDir(), "test.txt")
	uploadDir := filepath.Join(cwd, "upload/")
	templatesDir := filepath.Join(cwd, "templates/")

	if err := createTestFile(testFilePath); err != nil {
		t.Fatalf("cannot create test file: %s", err)
	}
	defer os.Remove(testFilePath)

	file, err := os.Open(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok!")
	}))
	defer ts.Close()

	uploadHandler := &UploadHandler{
		HostAddr:     ts.URL,
		UploadDir:    uploadDir,
		TemplatesDir: templatesDir,
	}

	uploadHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `<h1>Your file is uploaded</h1>`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFilesHandler_ServeHTTP(t *testing.T) {
	cwd := filepath.Join(getCurrentPath(), "..", "..")

	req, err := http.NewRequest("GET", "/?ext=txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := &FilesHandler{
		UploadHandler{
			HostAddr:     "127.0.0.1:80",
			UploadDir:    filepath.Join(cwd, "upload/"),
			TemplatesDir: filepath.Join(cwd, "templates/"),
		},
	}

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `<td><a href="files/test.txt">test</a></td>`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
