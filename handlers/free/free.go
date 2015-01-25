package free

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"

	alog "github.com/cenkalti/log"
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
	cmd := exec.Command("du", "-bcs", h.rootPath+"storage")

	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Run()

	readerFromOut := bytes.NewReader(out.Bytes())
	readerFromErr := bytes.NewReader(stderr.Bytes())

	rM := io.MultiReader(readerFromOut, readerFromErr)

	commandOutputComplete, err := ioutil.ReadAll(rM)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sizeInB, err := extractSizeFromString(string(commandOutputComplete))
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	procede := true
	if sizeInB > 50000000000 {
		procede = false
	}

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

func extractSizeFromString(source string) (int64, error) {
	re := regexp.MustCompile("^([0-9]+)")

	submatches := re.FindStringSubmatch(source)
	if len(submatches) < 2 {
		return 0, errors.New("Size not found")
	}

	i, err := strconv.ParseInt(submatches[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}
