package check

import (
	"encoding/json"
	"net/http"

	alog "github.com/cenkalti/log"
	"github.com/gophergala/videq/handlers/session"
	"github.com/gophergala/videq/janitor"
	"github.com/gophergala/videq/mediatools"
)

type VideoDimension struct {
	Height string
	Width  string
}

type VideoInfo struct {
	Procede          bool
	Err              string
	OutputDimensions []VideoDimension
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

	ok, mediaInfo, err := janitor.PossibleToEncode(sid)
	if err != nil {
		returnValue := VideoInfo{
			Procede:          ok,
			Err:              err.Error(),
			OutputDimensions: make([]VideoDimension, 0),
			OriginalInfo:     mediaInfo}
		js, err := json.Marshal(returnValue)
		if err != nil {
			h.log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}

	outputDm := make([]VideoDimension, 0)
	outputDm = append(outputDm, VideoDimension{"100", "200"})
	outputDm = append(outputDm, VideoDimension{"200", "400"})
	outputDm = append(outputDm, VideoDimension{"300", "600"})
	outputDm = append(outputDm, VideoDimension{"400", "800"})

	returnValue := VideoInfo{
		Procede:          ok,
		Err:              "",
		OutputDimensions: outputDm,
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
