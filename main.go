package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gophergala/videq/handlers/gzip"
	"github.com/gophergala/videq/handlers/home"
	"github.com/gophergala/videq/handlers/static"
)

const ROOT_PATH = "./"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := createStorage()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	staticHandler := static.NewHandler(ROOT_PATH)
	gzipStaticHandler := gzip.NewHandler(staticHandler)
	http.Handle("/resources/", gzipStaticHandler)

	homeHandler := home.NewHandler(ROOT_PATH)
	http.Handle("/", homeHandler)

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
