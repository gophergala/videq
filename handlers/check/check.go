package check

import (
	"encoding/json"
	"net/http"

	alog "github.com/cenkalti/log"
	"github.com/gophergala/videq/handlers/session"
	"github.com/gophergala/videq/janitor"
)

type VideoDimension struct {
	Height string
	Width  string
}

type VideoInfo struct {
	Err              string
	OutputDimensions []VideoDimension
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

	// contact media check

	outputDm := make([]VideoDimension, 0)
	outputDm = append(outputDm, VideoDimension{"100", "200"})
	outputDm = append(outputDm, VideoDimension{"200", "400"})
	outputDm = append(outputDm, VideoDimension{"300", "600"})
	outputDm = append(outputDm, VideoDimension{"400", "800"})

	returnValue := VideoInfo{
		Err:              "",
		OutputDimensions: outputDm}

	js, err := json.Marshal(returnValue)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	// stao - prebacit da se sejva u SID folder usera
}
