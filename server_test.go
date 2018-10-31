package gcs_proxy

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/http"
)

type StubRepository struct {
	getObjects func(path string) []string
}

func (s StubRepository) GetObjects(path string) []string {
	if s.getObjects == nil {
		return []string{}
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
	files := []string{
		"file1",
		"file2",
		"file3",
	}
	repository := StubRepository{}
	repository.getObjects = func(path string) []string {
		return files
	}

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(NewServer(repository).Handler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Contains(t, string(rr.Body.Bytes()), files[0])
	assert.Contains(t, string(rr.Body.Bytes()), files[1])
	assert.Contains(t, string(rr.Body.Bytes()), files[2])
}

