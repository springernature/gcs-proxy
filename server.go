package gcs_proxy

import (
	"net/http"
	"fmt"
)

type Server struct {
	repository ObjectRepository
}

func (s Server) Handler(w http.ResponseWriter, r *http.Request) {
	objects := s.repository.GetObjects("")

	var blah string
	for _, object := range objects {
		blah += fmt.Sprintf("%s\n", object)
	}
	w.Write([]byte(blah))
}

func NewServer(repository ObjectRepository) Server {
	return Server{
		repository: repository,
	}
}
