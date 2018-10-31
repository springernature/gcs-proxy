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
		<a href="../">../</a></br>
		{{range .}}
			<a href="/{{.Path}}">{{.Name}}</a></br>
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
	path := r.URL.Path
	if path == "/favicon.ico" {
		return
	}

		if path != "/" {
		isFile, err := s.repository.IsFile(path)
				if err != nil {
			handleError(err, w)
			return
		}
		if isFile {
						object, err := s.repository.GetObject(path)
			if err != nil {
				handleError(err, w)
				return
			}
			w.Write(object)
			return
		}
	}

		objects, err := s.repository.GetObjects(path)
	if err != nil {
		handleError(err, w)
		return
	}
	err = s.renderTemplate(objects, w)
	if err != nil {
		handleError(err, w)
		return
	}
}

func handleError(err error, w http.ResponseWriter) {
	log.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func NewServer(repository ObjectRepository) Server {
	return Server{
		repository: repository,
	}
}
