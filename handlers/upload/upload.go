package upload

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Handler struct {
	rootPath                  string
	numOfCompleteFileChBuffer int
	numOfCompleteFileWorkers  int
	completedFilesCh          chan string
}

func NewHandler(rootPath string, numOfCompleteFileChBuffer int, numOfCompleteFileWorkers int) *Handler {
	h := new(Handler)
	h.rootPath = rootPath
	h.numOfCompleteFileChBuffer = numOfCompleteFileChBuffer
	h.numOfCompleteFileWorkers = numOfCompleteFileWorkers

	h.completedFilesCh = make(chan string, h.numOfCompleteFileChBuffer)

	for i := 0; i < h.numOfCompleteFileWorkers; i++ {
		go assembleFile(h)
	}

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.streamUpload(w, r)
	} else if r.Method == "GET" {
		h.checkPart(w, r)
	} else {
		http.Error(w, "Method not suported", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) checkPart(w http.ResponseWriter, r *http.Request) {
	chunkDirPath := h.rootPath + "/storage/.upload/" + r.FormValue("flowFilename") + "/" + r.FormValue("flowChunkNumber")

	stat, err := os.Stat(chunkDirPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	expectedChunkSize, err := strconv.Atoi(r.URL.Query().Get("flowCurrentChunkSize"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if int(stat.Size()) != expectedChunkSize {
		http.Error(w, "Chunk size check failed", http.StatusPartialContent)
		return
	}
}

func (h *Handler) streamUpload(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	reader, err := r.MultipartReader()
	// Part 1: Chunk Number
	// Part 4: Total Size (bytes)
	// Part 6: File Name
	// Part 8: Total Chunks
	// Part 9: Chunk Data
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	part, err := reader.NextPart() // 1
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(buf, part)
	chunkNo := buf.String()
	buf.Reset()

	for i := 0; i < 3; i++ { // 2 3 4
		// move through unused parts
		part, err = reader.NextPart()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	io.Copy(buf, part)
	fileName := buf.String()
	buf.Reset()

	for i := 0; i < 3; i++ { // 7 8 9
		// move through unused parts
		part, err = reader.NextPart()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	chunkDirPath := h.rootPath + "/storage/.upload/" + fileName
	err = os.MkdirAll(chunkDirPath, 02750)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(chunkDirPath + "/" + chunkNo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, part)

	fileInfos, err := ioutil.ReadDir(chunkDirPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentSize := totalSize(fileInfos)
	flowTotalSizeInt64, err := strconv.ParseInt(flowTotalSize, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if flowTotalSizeInt64 == currentSize {
		h.completedFilesCh <- chunkDirPath
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
	for path := range h.completedFilesCh {

		fileInfos, err := ioutil.ReadDir(path)
		if err != nil {
			log.Print(err)
			return
		}

		// create final file to write to
		dst, err := os.Create(h.rootPath + "/storage/datastore/" + strings.Split(path, "/")[4])
		if err != nil {
			log.Print(err)
			return
		}
		defer dst.Close()

		sort.Sort(ByChunk(fileInfos))
		for _, fs := range fileInfos {
			src, err := os.Open(path + "/" + fs.Name())
			if err != nil {
				log.Print(err)
				return
			}
			defer src.Close()
			io.Copy(dst, src)
		}
		os.RemoveAll(path)
	}
}
