package static

import (
	"net/http"
	"os"
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
	requestedPath := r.URL.Path[len("/"):]

	file, err := os.Open(h.rootPath + requestedPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, requestedPath, stat.ModTime(), file)
}
