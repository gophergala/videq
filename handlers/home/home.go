package home

import (
	"html/template"
	"net/http"
)

type Handler struct {
	rootPath string
}

func NewHandler(rootPath string) *Handler {
	h := new(Handler)
	h.rootPath = rootPath
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("home.html").ParseFiles(h.rootPath + "templates/home.html"))

	p := make(map[string]interface{})

	tpl.Execute(w, p)
}
