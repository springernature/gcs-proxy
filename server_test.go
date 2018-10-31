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
	getObject  func(path string) ([]byte, error)
	isFile     func(path string) bool
}

func (s StubRepository) GetObject(path string) ([]byte, error) {
	if s.getObject == nil {
		return []byte(""), nil
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
		{Path: "imADirectory/directoryFile1", Name: "directoryFile1"},
		{Path: "imADirectory/directoryFile2", Name: "directoryFile1"},
		{Path: "imADirectory/directoryFile3", Name: "directoryFile1"},
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
	assert.Contains(t, string(rr.Body.Bytes()), directoryFiles[0].Path)
	assert.Contains(t, string(rr.Body.Bytes()), directoryFiles[1].Path)
	assert.Contains(t, string(rr.Body.Bytes()), directoryFiles[2].Path)
}

func TestItReturnsTheContentOfTheFile(t *testing.T) {
	requestPath := "/file1"
	file1Content := "imAFile"

	repository := StubRepository{}
	repository.isFile = func(path string) bool {
		return path == requestPath
	}
	repository.getObject = func(path string) ([]byte, error) {
		if path == requestPath {
			return []byte(file1Content), nil
		}
		t.Fatal("Should not get here")
		return []byte(""), nil
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

func TestItReturnsAnyErrorFromGetObject(t *testing.T) {
	expectedError := errors.New("nooooes")
	filePath := "path/to/file"
	repository := StubRepository{}

	repository.isFile = func(path string) bool {
		if path == filePath {
			return true
		}
		t.Error("Should not get here")
		return false
	}

	repository.getObject = func(path string) ([]byte, error) {
		if path == filePath {
			return nil, expectedError
		}
		t.Error("Should not get here")
		return nil, nil

	}

	req, _ := http.NewRequest("GET", filePath, nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(NewServer(repository).Handler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Contains(t, string(rr.Body.Bytes()), expectedError.Error())

}

//Todo test getObject error case
