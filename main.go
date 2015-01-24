package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gophergala/videq/handlers/gzip"
	"github.com/gophergala/videq/handlers/home"
	"github.com/gophergala/videq/handlers/static"
	"github.com/gophergala/videq/handlers/upload"
)

const ROOT_PATH = "./"
const NUM_OF_UPLOAD_WORKERS = 10
const NUM_OF_UPLOAD_BUFFER = 100

var completedFiles = make(chan string, 100)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := createStorage()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	webPort := flag.Int("web", 0, "Start web server on given port")
	flag.Parse()

	if *webPort > 0 {
		webServer()
		return
	}

	log.Println("Non server mode active")
}

func webServer() {
	staticHandler := static.NewHandler(ROOT_PATH)
	gzipStaticHandler := gzip.NewHandler(staticHandler)
	http.Handle("/resources/", gzipStaticHandler)

	homeHandler := home.NewHandler(ROOT_PATH)
	http.Handle("/", homeHandler)

	uploadHandler := upload.NewHandler(ROOT_PATH, NUM_OF_UPLOAD_BUFFER, NUM_OF_UPLOAD_WORKERS)
	http.Handle("/upload/", uploadHandler)

	log.Println("Server started on port 8094")
	log.Print(http.ListenAndServe(":8094", nil))
}

func createStorage() error {
	paths := []string{
		ROOT_PATH + "storage/datastore",
		ROOT_PATH + "storage/.upload"}

	for _, path := range paths {
		err := os.MkdirAll(path, 02750)
		if err != nil {
			return err
		}
	}

	return nil
}
