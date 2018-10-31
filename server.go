package gcs_proxy

import (
	"net/http"
	"html/template"
	"github.com/JFrogDev/jfrog-cli-go/jfrog-client/utils/log"
)

type Server struct {
	repository ObjectRepository
}

func (s Server) renderTemplate(objects []Object, w http.ResponseWriter) (err error) {
	t := template.New("template")
	template, err := t.Parse(`
<html>
	<body>
		{{range .}}
			<a href="{{.Path}}">{{.Name}}</a>
		{{end}}
	</body>
</html>

`)
	if err != nil {
		return
	}

	return template.Execute(w, objects)
}

func (s Server) writeFile(path string, w http.ResponseWriter) {
	object := s.repository.GetObject(path)
	w.Write(object)
}

func (s Server) Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path != "/" && s.repository.IsFile(path) {
		w.Write(s.repository.GetObject(path))
		return
	}

	objects, err := s.repository.GetObjects(path)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
	err = s.renderTemplate(objects, w)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func NewServer(repository ObjectRepository) Server {
	return Server{
		repository: repository,
	}
}
