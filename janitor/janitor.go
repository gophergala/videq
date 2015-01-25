package janitor

import (
	"database/sql"
	"os"
)

var StorageIncomplete string
var StorageComplete string
var DbConn *sql.DB

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
