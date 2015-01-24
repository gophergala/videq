package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gophergala/videq/handlers/gzip"
	"github.com/gophergala/videq/handlers/home"
	"github.com/gophergala/videq/handlers/session"
	"github.com/gophergala/videq/handlers/static"
	"github.com/gophergala/videq/handlers/upload"
)

const ROOT_PATH = "./"
const NUM_OF_MERGE_WORKERS = 10
const NUM_OF_MERGE_BUFFER = 100
const DSN = "root:m11@/videq"

var completedFiles = make(chan string, 100)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := createStorage()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	webPort := flag.String("web", "", "Start web server on given port")
	flag.Parse()

	if *webPort != "" {
		webServer(*webPort)
		return
	}

	log.Println("Non server mode active")
}

func webServer(port string) {
	staticHandler := static.NewHandler(ROOT_PATH)
	gzipStaticHandler := gzip.NewHandler(staticHandler)
	http.Handle("/resources/", gzipStaticHandler)

	homeHandler := home.NewHandler(ROOT_PATH)
	homeSidHandler := session.NewHandler(homeHandler, DSN)
	http.Handle("/", homeSidHandler)

	uploadHandler := upload.NewHandler(ROOT_PATH, NUM_OF_MERGE_BUFFER, NUM_OF_MERGE_WORKERS)
	http.Handle("/upload/", uploadHandler)

	log.Printf("Server started on port %v", port)
	log.Print(http.ListenAndServe(":"+port, nil))
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
