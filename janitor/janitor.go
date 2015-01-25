package janitor

import (
	"database/sql"
	"os"

	"github.com/gophergala/videq/mediatools"
)

var StorageIncomplete string
var StorageComplete string
var DbConn *sql.DB
var log alog.Logger

var cleanUploadFolderCh chan string

func Init(db *sql.DB, sc, si string, l alog.Logger) {
	DbConn = db
	StorageComplete = sc
	StorageIncomplete = si
	log = l

	cleanUploadFolderCh = make(chan string, 100)
	for i := 0; i < 10; i++ {
		go clenupIncompleteFolder(cleanUploadFolderCh)
	}
}

// check if current user is uploading a file?
func HasFileInUpload(sid string) (bool, error) {
	firstPartFilename := StorageIncomplete + "/" + sid + "/1"

	_, err := os.Stat(firstPartFilename)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func RecordFilename(sid, filename string) error {
	_, err := DbConn.Exec("INSERT INTO file (sid, filename, start_ts) VALUES (?, ?, UNIX_TIMESTAMP())", sid, filename)
	if err != nil {
		return err
	}
	return nil
}

func PossibleToEncode(sid string) (bool, mediatools.MediaFileInfo) {
	mt := mediatools.NewMediaInfo(log)

	userFolder := StorageIncomplete + sid

	ok, minfob, err := mt.CheckMedia(userFolder + "/1")
	if err != nil {
		log.Error(err)
		cleanUploadFolderCh <- userFolder
		return false, minfob
	}

	return ok, minfob
}

func clenupIncompleteFolder(pathCh <-chan string) {
	for path := range pathCh {
		os.RemoveAll(path)
	}
}
