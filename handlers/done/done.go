package done

import (
	"database/sql"
	"encoding/json"
	"net/http"

	alog "github.com/cenkalti/log"
	"github.com/gophergala/videq/handlers/session"
	"github.com/gophergala/videq/janitor"
)

type Handler struct {
	rootPath string
	log      alog.Logger
	db       *sql.DB
}

func NewHandler(log alog.Logger, rootPath string, db *sql.DB) *Handler {
	h := new(Handler)
	h.log = log
	h.rootPath = rootPath
	h.db = db
	return h
}

type downloadData struct {
	Procede         bool
	Err             string
	First_frame_jpg string `json:"first_frame_jpg"`
	Mp4_link        string `json:"mp4_link"`
	Webm_link       string `json:"webm_link"`
	Ogv_link        string `json:"ogv_link"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid, err := session.Sid(r)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dd := &downloadData{}

	var success sql.NullInt64
	var encodingErr sql.NullString

	err = h.db.QueryRow("SELECT success, encode_error FROM file WHERE sid=? ", sid).Scan(&success, &encodingErr)
	switch {
	case err == sql.ErrNoRows:
		dd.Procede = false
		dd.Err = "No encoding job found"
		h.log.Errorf("No encoding job found for sid=%v", sid)
		break

	case err != nil:
		h.log.Error(err)
		dd.Procede = false
		dd.Err = err.Error()
		break

	default:
		if success.Int64 > 0 {
			dd.Procede = true
			dd.Err = encodingErr.String
			dd.First_frame_jpg = "/download/encoded.jpg"
			dd.Mp4_link = "/download/encoded.mp4"
			dd.Ogv_link = "/download/encoded.ogg"
			dd.Webm_link = "/download/encoded.webm"
		} else if len(encodingErr.String) > 0 {
			dd.Procede = false
			dd.Err = encodingErr.String
			janitor.CleanupUser(sid)
		} else {
			dd.Procede = false
		}
	}

	js, err := json.Marshal(dd)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
