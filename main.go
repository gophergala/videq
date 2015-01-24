package main

import (
	"flag"
	//	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/cenkalti/log"
	"github.com/gophergala/videq/handlers/gzip"
	"github.com/gophergala/videq/handlers/home"
	"github.com/gophergala/videq/handlers/session"
	"github.com/gophergala/videq/handlers/static"
	"github.com/gophergala/videq/handlers/upload"
	"github.com/gophergala/videq/janitor"
	"github.com/gophergala/videq/mediatools"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const ROOT_PATH = "./"
const NUM_OF_MERGE_WORKERS = 10
const NUM_OF_MERGE_BUFFER = 100
const DSN = "root:m11@/videq"

var db *Database

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	InitLogger()
	LoadConfig()

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

	minfo, err := mediatools.GetMediaInfo()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("%#v", minfo)
}

func webServer(db *Database, port string) {
	staticHandler := static.NewHandler(ROOT_PATH)
	gzipStaticHandler := gzip.NewHandler(staticHandler)
	http.Handle("/resources/", gzipStaticHandler)

	homeHandler := home.NewHandler(ROOT_PATH)
	homeSidHandler := session.NewHandler(log, db.conn, homeHandler, DSN)
	http.Handle("/", homeSidHandler)

	uploadHandler := upload.NewHandler(log, ROOT_PATH, NUM_OF_MERGE_BUFFER, NUM_OF_MERGE_WORKERS)
	http.Handle("/upload/", uploadHandler)

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

type Database struct {
	conn *sql.DB // global variable to share it between main and short lived functions (and eg. the HTTP handler)
}

type DbConfig struct {
	DbHost string
	DbName string
	DbUser string
	DbPass string
	Debug  bool
}

//func OpenDB(host, name, user, pass string) *Database {
func (d *Database) OpenDB(cfg DbConfig) (error, *Database) {

	dba, err := sql.Open("mysql", cfg.DbUser+":"+cfg.DbPass+"@tcp("+cfg.DbHost+":3306)/"+cfg.DbName+"?charset=utf8mb4,utf8")
	if err != nil {
		// VAZNO: NIKAD SE NE OKINE!!!! (al svejedno treba provjeravat ... XXX TODO doh)
		//log.Debug(err)
		//		log.Debug("debug %s", Password("secret"))
		return err, nil
	}
	// prebacio u fju iznad koji poziva OpenDB. zasto? da mi se ne zatvori kad izadjem iz ove fje?
	//	defer db.Close()

	err = dba.Ping() // zato se cesto koristi ping
	if err != nil {
		//log.Fatal(err)
		return err, nil
	}
	dba.SetMaxIdleConns(100)
	dba.SetMaxOpenConns(200)

	return nil, &Database{conn: dba}
}

func (d *Database) CloseDB() {
	d.conn.Close()
}
