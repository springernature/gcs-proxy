package gcs_proxy

import (
	"net/http"
	"html/template"
)

type Server struct {
	repository ObjectRepository
}

func (s Server) renderTemplate(objects []string, w http.ResponseWriter) (err error) {
	t := template.New("template")
	template, err := t.Parse(`
<html>
	<body>
		{{range .}}
			<a href="{{.}}">{{.}}</a>
		{{end}}
	</body>
</html>

`)
	if err != nil {
		return
	}

	return template.Execute(w, objects)
}

func (s Server) Handler(w http.ResponseWriter, r *http.Request) {
	objects := s.repository.GetObjects("")
	err := s.renderTemplate(objects, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewServer(repository ObjectRepository) Server {
	return Server{
		repository: repository,
	}
}
