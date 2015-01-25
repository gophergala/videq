package download

import (
	"net/http"
	"os"
	"strings"

	alog "github.com/cenkalti/log"
	"github.com/gophergala/videq/handlers/session"
)

type Handler struct {
	rootPath string
	log      alog.Logger
}

func NewHandler(log alog.Logger, rootPath string) *Handler {
	h := new(Handler)
	h.rootPath = rootPath
	h.log = log
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid, err := session.Sid(r)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	urlParts := strings.Split(r.URL.Path, "/")
	filename := urlParts[len(urlParts)-1]

	requestedPath := h.rootPath + "storage/datastore/" + sid + "/" + filename

	file, err := os.Open(requestedPath)
	if err != nil {
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	stat, err := file.Stat()
	if err != nil {
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Debugf("Download %v/%v", sid, filename)

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeContent(w, r, requestedPath, stat.ModTime(), file)
}
