package janitor

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	alog "github.com/cenkalti/log"
	"github.com/gophergala/videq/mediatools"
)

var Storage string
var StorageIncomplete string
var StorageComplete string
var DbConn *sql.DB
var log alog.Logger

var cleanUploadFolderCh chan string
var encodePathCh chan string
var allowedToUlCh chan bool

func Init(db *sql.DB, s, sc, si string, l alog.Logger) {
	DbConn = db
	Storage = s
	StorageComplete = sc
	StorageIncomplete = si
	log = l

	cleanUploadFolderCh = make(chan string, 100)
	for i := 0; i < 10; i++ {
		go cleanupIncompleteFolderWorker(cleanUploadFolderCh)
	}

	encodePathCh = make(chan string, 1000000)
	for i := 0; i < 3; i++ {
		go encodeWorker(encodePathCh)
	}

	allowedToUlCh = make(chan bool)
	go freeSpaceCheckWorker()
}

func CleanupUser(sid string) error {
	os.RemoveAll(StorageComplete + sid)
	os.RemoveAll(StorageIncomplete + sid)

	_, err := DbConn.Exec("DELETE FROM file WHERE sid=?", sid)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func IsAllowedToUpload() bool {
	return <-allowedToUlCh
}

// check if current user is uploading a file?
func HasFileInUpload(sid string) (bool, error) {
	firstPartFilename := StorageIncomplete + sid + "/1"

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
	_, err := DbConn.Exec("REPLACE INTO file (sid, filename, start_ts) VALUES (?, ?, UNIX_TIMESTAMP())", sid, filename)
	if err != nil {
		return err
	}
	return nil
}

func PossibleToEncode(sid string) (bool, mediatools.MediaFileInfo, map[string]mediatools.VideoResolution, error) {
	mt := mediatools.NewMediaInfo(log)

	userFolder := StorageIncomplete + sid

	ok, minfob, resolutions, err := mt.CheckMedia(userFolder + "/1")
	if err != nil {
		log.Error(err)
		cleanUploadFolderCh <- userFolder
		return false, minfob, resolutions, err
	}

	return ok, minfob, resolutions, nil
}

func PushToEncode(path string) {
	sid := strings.Split(path, "/")[2] // todo - make batter

	_, err := DbConn.Exec("UPDATE file SET path_of_original=?, added_to_encode_queue_ts=UNIX_TIMESTAMP() WHERE sid=?", path, sid)
	if err != nil {
		log.Error(err)
		// todo whole cleanup
		return
	}

	encodePathCh <- path
}

func encodeWorker(pathCh <-chan string) {
	for path := range pathCh {
		sid := strings.Split(path, "/")[2] // todo - make batter

		_, err := DbConn.Exec("UPDATE file SET encode_start_ts=UNIX_TIMESTAMP() WHERE sid=?", sid)
		if err != nil {
			log.Error(err)
			// todo whole cleanup
			continue
		}

		pathSpl := strings.Split(path, "/")
		filePath := "./" + strings.Join(pathSpl[0:len(pathSpl)-1], "/") + "/"
		fileName := pathSpl[len(pathSpl)-1]
		log.Debugln(path, filePath, fileName)

		mt := mediatools.NewMediaInfo(log)
		err = mt.EncodeVideoFile(filePath, fileName)
		if err != nil {
			log.Error(err)
		}

		encodeEnded(sid, err)
	}
}

func encodeEnded(sid string, err error) {
	errorTxt := ""
	success := 1

	if err != nil {
		errorTxt = err.Error()
		success = 0
	}

	_, err = DbConn.Exec("UPDATE file SET encode_end_ts=UNIX_TIMESTAMP(), encode_error=?, success=? WHERE sid=?", errorTxt, success, sid)
	if err != nil {
		log.Error(err)
		return
	}
}

func cleanupIncompleteFolderWorker(pathCh <-chan string) {
	for path := range pathCh {
		os.RemoveAll(path)
	}
}

func freeSpaceCheckWorker() {
	isFree := true

	ct := time.Tick(2 * time.Minute)

	for {
		select {
		case allowedToUlCh <- isFree:
			break
		case <-ct:
			cmd := exec.Command("du", "-bcs", Storage)

			var out bytes.Buffer
			cmd.Stdout = &out

			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			cmd.Run()

			readerFromOut := bytes.NewReader(out.Bytes())
			readerFromErr := bytes.NewReader(stderr.Bytes())

			rM := io.MultiReader(readerFromOut, readerFromErr)

			commandOutputComplete, err := ioutil.ReadAll(rM)
			if err != nil {
				log.Error(err)
				break
			}

			sizeInB, err := extractSizeFromString(string(commandOutputComplete))
			if err != nil {
				log.Error(err)
				break
			}

			isFree = true
			if sizeInB > 50000000000 {
				isFree = false
			}
			break
		}
	}
}

func extractSizeFromString(source string) (int64, error) {
	re := regexp.MustCompile("^([0-9]+)")

	submatches := re.FindStringSubmatch(source)
	if len(submatches) < 2 {
		return 0, errors.New("Size not found")
	}

	i, err := strconv.ParseInt(submatches[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}
