package gcs_proxy

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/http"
	"errors"
)

type StubRepository struct {
	getObjects func(path string) ([]Object, error)
	getObject  func(path string) []byte
	isFile     func(path string) bool
}

func (s StubRepository) GetObject(path string) []byte {
	if s.getObject == nil {
		return []byte("")
	}
	return s.getObject(path)

}

func (s StubRepository) IsFile(path string) bool {
	if s.isFile == nil {
		return false
	}

	return s.isFile(path)
}

func (s StubRepository) GetObjects(path string) ([]Object, error) {
	if s.getObjects == nil {
		return []Object{}, nil
	}
	return s.getObjects(path)
}

func TestItRespondsWithA200OnSlash(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NewServer(StubRepository{}).Handler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}

func TestItRendersFiles(t *testing.T) {
	files := []Object{
		{Name: "file1"},
		{Name: "file2"},
		{Name: "file3"},
	}
	repository := StubRepository{}
	repository.getObjects = func(path string) ([]Object, error) {
		return files, nil
	}

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(NewServer(repository).Handler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Contains(t, string(rr.Body.Bytes()), files[0].Name)
	assert.Contains(t, string(rr.Body.Bytes()), files[1].Name)
	assert.Contains(t, string(rr.Body.Bytes()), files[2].Name)
}

func TestItRendersFilesWhenTraversing(t *testing.T) {
	slashPath := "/"
	slashFiles := []Object{
		{Name: "file1"},
		{Name: "file2"},
		{Name: "imADirectory"},
	}

	directoryPath := "/imADirectory"
	directoryFiles := []Object{
		{Name: "directoryFile1"},
		{Name: "directoryFile2"},
		{Name: "directoryFile3"},
	}
	repository := StubRepository{}
	repository.getObjects = func(path string) ([]Object, error) {
		if path == slashPath {
			return slashFiles, nil
		}
		if path == directoryPath {
			return directoryFiles, nil
		}
		t.Error("Should not get here")
		return []Object{}, nil
	}

	req, _ := http.NewRequest("GET", slashPath, nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(NewServer(repository).Handler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Contains(t, string(rr.Body.Bytes()), slashFiles[0].Name)
	assert.Contains(t, string(rr.Body.Bytes()), slashFiles[1].Name)
	assert.Contains(t, string(rr.Body.Bytes()), slashFiles[2].Name)

	req, _ = http.NewRequest("GET", directoryPath, nil)
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(NewServer(repository).Handler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Contains(t, string(rr.Body.Bytes()), directoryFiles[0].Name)
	assert.Contains(t, string(rr.Body.Bytes()), directoryFiles[1].Name)
	assert.Contains(t, string(rr.Body.Bytes()), directoryFiles[2].Name)
}

func TestItReturnsTheContentOfTheFile(t *testing.T) {
	requestPath := "/file1"
	file1Content := "imAFile"

	repository := StubRepository{}
	repository.isFile = func(path string) bool {
		return path == requestPath
	}
	repository.getObject = func(path string) []byte {
		if path == requestPath {
			return []byte(file1Content)
		}
		t.Fatal("Should not get here")
		return []byte("")
	}

	req, _ := http.NewRequest("GET", requestPath, nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(NewServer(repository).Handler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, string(rr.Body.Bytes()), file1Content)
}

func TestItReturnsAnyErrorFromGetObjects(t *testing.T) {
	expectedError := errors.New("nooooes")

	repository := StubRepository{}
	repository.getObjects = func(path string) ([]Object, error) {
		return []Object{}, expectedError
	}

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(NewServer(repository).Handler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Contains(t, string(rr.Body.Bytes()), expectedError.Error())
}
