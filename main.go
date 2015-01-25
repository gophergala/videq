package main

import (
	"flag"
	"net/http"
	"os"
	"runtime"

	_ "github.com/cenkalti/log"
	"github.com/gophergala/videq/config"
	"github.com/gophergala/videq/handlers/check"
	"github.com/gophergala/videq/handlers/done"
	"github.com/gophergala/videq/handlers/download"
	"github.com/gophergala/videq/handlers/free"
	"github.com/gophergala/videq/handlers/gzip"
	"github.com/gophergala/videq/handlers/home"
	"github.com/gophergala/videq/handlers/restart"
	"github.com/gophergala/videq/handlers/session"
	"github.com/gophergala/videq/handlers/static"
	"github.com/gophergala/videq/handlers/upload"
	"github.com/gophergala/videq/janitor"
	//"github.com/gophergala/videq/mediatools"
)

const ROOT_PATH = "./"
const NUM_OF_MERGE_WORKERS = 10
const NUM_OF_MERGE_BUFFER = 100

var cfg config.Config
var db *Database

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	InitLogger()
	config.LoadConfig(log, &cfg)
	checkExecutables()

	err := createStorage()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	err, db = db.OpenDB(DbConfig{cfg.DB.HOST, cfg.DB.NAME, cfg.DB.USER, cfg.DB.PASS, cfg.DB.DEBUG})
	if db == nil {
		log.Fatal("Error, cannot connect to db, db.OpenDB ", err)
	}
	defer db.CloseDB()

	janitor.Init(db.conn, ROOT_PATH+"storage", ROOT_PATH+"storage/datastore/", ROOT_PATH+"storage/.upload/", log)

	webPort := flag.String("web", "", "Start web server on given port")
	flag.Parse()

	if *webPort != "" {
		webServer(db, *webPort)
		return
	}

	log.Infoln("Non server mode active")

	// -------------------------------------------------------
	// beware, hard testing bewlow
	// -------------------------------------------------------

	// mt := mediatools.NewMediaInfo(log)
	// _ = mt

	// testiranje mediainfo toola
	// minfo, err := mt.GetMediaInfo("_test/master_1080.mp4")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Infof("%#v", minfo)

	// testiranje checkera
	// ok, minfob, res, err := mt.CheckMedia("_test/master_1080.mp4") //  "test.psd" "r2w_1080p.mov" _test/master_1080.mp4" "videq_sw.mp4"
	// if err != nil {
	// 	log.Error(err)
	// }
	// log.Infof("%#v", ok)
	// log.Infof("%#v", minfob)
	// log.Infof("%#v", res)
	// log.Infof("%#v", err)

	//err = mt.EncodeVideoFile("_test/", "videq_sw.mp4") // "master_1080.mp4"   "r2w_1080p.mov"

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

	doneHandler := done.NewHandler(log, ROOT_PATH, db.conn)
	http.Handle("/done/", doneHandler)

	downloadHandler := download.NewHandler(log, ROOT_PATH)
	http.Handle("/download/", downloadHandler)

	restartHandler := restart.NewHandler(log)
	http.Handle("/restart/", restartHandler)

	freeHandler := free.NewHandler(log, ROOT_PATH)
	http.Handle("/free/", freeHandler)

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
