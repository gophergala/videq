package restart

import (
	"net/http"

	alog "github.com/cenkalti/log"
	"github.com/gophergala/videq/handlers/session"
	"github.com/gophergala/videq/janitor"
)

type Handler struct {
	log alog.Logger
}

func NewHandler(log alog.Logger) *Handler {
	h := new(Handler)
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

	err = janitor.CleanupUser(sid)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Debugf("User %v cleanup on damand", sid)
}
