package main

import (
	"flag"
	//	"log"
	"net/http"
	"os"
	"runtime"

	_ "github.com/cenkalti/log"
	"github.com/gophergala/videq/handlers/check"
	"github.com/gophergala/videq/handlers/gzip"
	"github.com/gophergala/videq/handlers/home"
	"github.com/gophergala/videq/handlers/session"
	"github.com/gophergala/videq/handlers/static"
	"github.com/gophergala/videq/handlers/upload"
	"github.com/gophergala/videq/janitor"
	"github.com/gophergala/videq/mediatools"
)

const ROOT_PATH = "./"
const NUM_OF_MERGE_WORKERS = 10
const NUM_OF_MERGE_BUFFER = 100

var db *Database

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	InitLogger()
	LoadConfig()
	checkExecutables()

	err := createStorage()
	if err != nil {
		log.Fatal(err)
	}

	janitor.StorageComplete = ROOT_PATH + "storage/datastore/"
	janitor.StorageIncomplete = ROOT_PATH + "storage/.upload/"
}

func main() {
	var err error
	err, db = db.OpenDB(DbConfig{config.DB.HOST, config.DB.NAME, config.DB.USER, config.DB.PASS, config.DB.DEBUG})
	if db == nil {
		log.Fatal("Error, cannot connect to db, db.OpenDB ", err)
	}
	defer db.CloseDB()

	webPort := flag.String("web", "", "Start web server on given port")
	flag.Parse()

	if *webPort != "" {
		webServer(db, *webPort)
		return
	}

	log.Infoln("Non server mode active")

	mt := mediatools.NewMediaInfo(log)

	// minfo, err := mt.GetMediaInfo("_test/master_1080.mp4")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Infof("%#v", minfo)

	ok, minfob, err := mt.CheckMedia("_test/videq_sw.mp4") // "_test/master_1080.mp4"
	if err != nil {
		log.Error(err)
	}
	log.Infof("%#v", ok)
	log.Infof("%#v", minfob)
	log.Infof("%#v", err)

}

func webServer(db *Database, port string) {
	staticHandler := static.NewHandler(ROOT_PATH)
	gzipStaticHandler := gzip.NewHandler(staticHandler)
	http.Handle("/resources/", gzipStaticHandler)

	homeHandler := home.NewHandler(ROOT_PATH)
	homeSidHandler := session.NewHandler(log, db.conn, homeHandler)
	http.Handle("/", homeSidHandler)

	uploadHandler := upload.NewHandler(log, ROOT_PATH, NUM_OF_MERGE_BUFFER, NUM_OF_MERGE_WORKERS)
	http.Handle("/upload/", uploadHandler)

	checkHandler := check.NewHandler(log, ROOT_PATH)
	http.Handle("/check/", checkHandler)

	log.Infof("Server started on port %v", port)
	log.Info(http.ListenAndServe(":"+port, nil))
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
