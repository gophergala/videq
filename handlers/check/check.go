package check

import (
	"encoding/json"
	"net/http"

	alog "github.com/cenkalti/log"
	"github.com/gophergala/videq/handlers/session"
	"github.com/gophergala/videq/janitor"
	"github.com/gophergala/videq/mediatools"
)

type VideoInfo struct {
	Procede          bool
	Err              string
	OutputDimensions map[string]mediatools.VideoResolution
	OriginalInfo     mediatools.MediaFileInfo
}

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

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid, err := session.Sid(r)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	hasFileInUpload, err := janitor.HasFileInUpload(sid)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if hasFileInUpload == false {
		h.log.Error("No upload to check for client " + sid)
		http.Error(w, "No file on server", http.StatusNotFound)
		return
	}

	ok, mediaInfo, res, err := janitor.PossibleToEncode(sid)
	errorString := ""
	if err != nil {
		errorString = err.Error()
	}
	returnValue := VideoInfo{
		Procede:          ok,
		Err:              errorString,
		OutputDimensions: res,
		OriginalInfo:     mediaInfo}

	js, err := json.Marshal(returnValue)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
