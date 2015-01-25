package upload

import (
	"bytes"
	"io"
	"io/ioutil"
	"path"
	//	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	alog "github.com/cenkalti/log"
	"github.com/gophergala/videq/handlers/session"
	"github.com/gophergala/videq/janitor"
)

type FileToAssemble struct {
	PathToParts      string
	OriginalFilename string
}

type Handler struct {
	rootPath                  string
	numOfCompleteFileChBuffer int
	numOfCompleteFileWorkers  int
	completedFilesCh          chan *FileToAssemble
	log                       alog.Logger
}

func NewHandler(log alog.Logger, rootPath string, numOfCompleteFileChBuffer int, numOfCompleteFileWorkers int) *Handler {
	h := new(Handler)
	h.log = log
	h.rootPath = rootPath
	h.numOfCompleteFileChBuffer = numOfCompleteFileChBuffer
	h.numOfCompleteFileWorkers = numOfCompleteFileWorkers

	h.completedFilesCh = make(chan *FileToAssemble, h.numOfCompleteFileChBuffer)

	for i := 0; i < h.numOfCompleteFileWorkers; i++ {
		go assembleFile(h)
	}

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isAllowedToUpload := janitor.IsAllowedToUpload()
	if isAllowedToUpload == false {
		http.Error(w, "No upload allowed at this time", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		h.streamUpload(w, r)
	} else if r.Method == "GET" {
		h.checkPart(w, r)
	} else {
		http.Error(w, "Method not suported", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) checkPart(w http.ResponseWriter, r *http.Request) {
	sid, err := session.Sid(r)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	chunkDirPath := h.rootPath + "storage/.upload/" + sid + "/" + r.FormValue("flowChunkNumber")

	stat, err := os.Stat(chunkDirPath)
	if err != nil {
		h.log.Debug(err)
		http.Error(w, "Internal server error", http.StatusNoContent)
		return
	}

	expectedChunkSize, err := strconv.Atoi(r.URL.Query().Get("flowCurrentChunkSize"))
	if err != nil {
		h.log.Debug(err)
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}
	if int(stat.Size()) != expectedChunkSize {
		h.log.Debug(err)
		http.Error(w, "Chunk size check failed", http.StatusPartialContent)
		return
	}
}

func (h *Handler) streamUpload(w http.ResponseWriter, r *http.Request) {
	sid, err := session.Sid(r)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	reader, err := r.MultipartReader()
	// Part 1: Chunk Number
	// Part 4: Total Size (bytes)
	// Part 6: File Name
	// Part 8: Total Chunks
	// Part 9: Chunk Data
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	part, err := reader.NextPart() // 1
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	io.Copy(buf, part)
	chunkNo := buf.String()
	buf.Reset()

	for i := 0; i < 3; i++ { // 2 3 4
		// move through unused parts
		part, err = reader.NextPart()
		if err != nil {
			h.log.Error(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	io.Copy(buf, part)
	flowTotalSize := buf.String()
	buf.Reset()

	for i := 0; i < 2; i++ { // 5 6
		// move through unused parts
		part, err = reader.NextPart()
		if err != nil {
			h.log.Error(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	io.Copy(buf, part)
	fileName := buf.String()
	buf.Reset()

	if chunkNo == "1" {
		err = janitor.RecordFilename(sid, fileName)
		if err != nil {
			h.log.Error(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	for i := 0; i < 3; i++ { // 7 8 9
		// move through unused parts
		part, err = reader.NextPart()
		if err != nil {
			h.log.Error(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	chunkDirPath := h.rootPath + "/storage/.upload/" + sid
	err = os.MkdirAll(chunkDirPath, 02750)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(chunkDirPath + "/" + chunkNo)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, part)

	fileInfos, err := ioutil.ReadDir(chunkDirPath)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	currentSize := totalSize(fileInfos)
	flowTotalSizeInt64, err := strconv.ParseInt(flowTotalSize, 10, 64)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if flowTotalSizeInt64 == currentSize {
		fta := &FileToAssemble{chunkDirPath, fileName}
		h.completedFilesCh <- fta
	}
}

func totalSize(fileInfos []os.FileInfo) int64 {
	var sum int64
	for _, fi := range fileInfos {
		sum += fi.Size()
	}
	return sum
}

type ByChunk []os.FileInfo

func (a ByChunk) Len() int      { return len(a) }
func (a ByChunk) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByChunk) Less(i, j int) bool {
	ai, _ := strconv.Atoi(a[i].Name())
	aj, _ := strconv.Atoi(a[j].Name())
	return ai < aj
}

func assembleFile(h *Handler) {
	for fta := range h.completedFilesCh {

		fileInfos, err := ioutil.ReadDir(fta.PathToParts)
		if err != nil {
			h.log.Error(err)
			return
		}

		sid := strings.Split(fta.PathToParts, "/")[4]

		targetDirPath := h.rootPath + "/storage/datastore/" + sid
		err = os.MkdirAll(targetDirPath, 02750)
		if err != nil {
			h.log.Error(err)
			return
		}

		partFilename := strings.Split(fta.OriginalFilename, ".")

		assabledFilePath := targetDirPath + "/original." + partFilename[len(partFilename)-1]

		// create final file to write to
		dst, err := os.Create(assabledFilePath)
		if err != nil {
			h.log.Error(err)
			return
		}
		defer dst.Close()

		sort.Sort(ByChunk(fileInfos))
		for _, fs := range fileInfos {
			func() {
				src, err := os.Open(fta.PathToParts + "/" + fs.Name())
				if err != nil {
					h.log.Error(err)
					return
				}
				defer src.Close()
				io.Copy(dst, src)
			}()
		}
		os.RemoveAll(fta.PathToParts)

		assabledFilePath = path.Clean(assabledFilePath)

		janitor.PushToEncode(assabledFilePath)
	}
}
