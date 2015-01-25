package janitor

import "os"

var StorageIncomplete string
var StorageComplete string

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
