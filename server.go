package gcs_proxy

import (
	"github.com/JFrogDev/jfrog-cli-go/jfrog-client/utils/log"
	"html/template"
	"mime"
	"net/http"
	"path"
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
	objPath := r.URL.Path
	if objPath == "/favicon.ico" {
		return
	}

	if objPath != "/" {
		isFile, err := s.repository.IsFile(objPath)
		if err != nil {
			handleError(err, w)
			return
		}
		if isFile {
			ext := path.Ext(objPath)
			contentType := mime.TypeByExtension(ext)

			if contentType != "" {
				w.Header().Set("Content-Type", contentType)
			}

			object, err := s.repository.GetObject(objPath)
			if err != nil {
				handleError(err, w)
				return
			}
			w.Write(object)

			return
		}
	}

	objects, err := s.repository.GetObjects(objPath)
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
