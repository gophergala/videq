package free

import (
	"encoding/json"
	"net/http"

	alog "github.com/cenkalti/log"
	"github.com/gophergala/videq/janitor"
)

type Handler struct {
	rootPath string
	log      alog.Logger
}

func NewHandler(log alog.Logger, rootPath string) *Handler {
	h := new(Handler)
	h.log = log
	h.rootPath = rootPath
	return h
}

type procedeValue struct {
	Procede bool
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	procede := janitor.IsAllowedToUpload()

	p := &procedeValue{}
	p.Procede = procede

	js, err := json.Marshal(p)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
