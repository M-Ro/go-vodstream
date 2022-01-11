package web

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed static/app.css static/app.js static/index.html.tmpl
var ContentFS embed.FS

// FsHandler handles static files (css, js, etc)
func FsHandler() http.Handler {
	sub, err := fs.Sub(ContentFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	return http.FileServer(http.FS(sub))
}

// IndexHandler returns the HTML template
func IndexHandler(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(ContentFS, "static/index.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	var path = req.URL.Path
	w.Header().Add("Content-Type", "text/html")
	t.Execute(w, struct {
		Title    string
		Response string
	}{Title: "VodStream", Response: path})
}
