package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	alog "github.com/cenkalti/log"
)

var log alog.Logger

// 	// ERROR - 2012-08-17 17:35:35 --> [VuduNonFatalException] D:\Inetpub\wwwvirtual\www.bonbon.hr\vudu_system\voodoo\URI.php: 170 [VuduNonFatalException] (The URI you submitted has disallowed characters. (g=js&121))  (404)
// 	// log_format := "%{color}[%{level:.4s}] %{time:15:04:05.000000} %{id:03x} %{shortfile} [%{longpkg}] %{longfunc} -> %{color:reset}%{message}"

func InitLogger() {
	processName := path.Base(os.Args[0])
	baseName := strings.Replace(processName, ".exe", "", -1)
	logFilename := fmt.Sprintf("%s.log", baseName)

	// Log levels (DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL)

	log = alog.NewLogger(processName)
	log.SetLevel(alog.DEBUG) // forward all messages to handler

	consoleLog := alog.NewWriterHandler(os.Stderr)
	consoleLog.SetFormatter(logFormatter{})
	consoleLog.SetLevel(alog.DEBUG)
	consoleLog.Colorize = true

	fileLog := alog.NewWriterHandler(logFile(logFilename))
	fileLog.SetLevel(alog.NOTICE)
	//log.SetHandler(fileLog)

	multi := alog.NewMultiHandler(consoleLog, fileLog)
	multi.SetFormatter(logFormatter{})

	log.SetHandler(multi)

}

type logFormatter struct{}

// %.4s limitira na 4
// %-4s padda lijevo 4 spacea

// Format outputs a message like "2014-02-28 18:15:57 [example] INFO     somethinfig happened"
func (f logFormatter) Format(rec *alog.Record) string {
	//	return fmt.Sprintf("%s %.4s [%s] %s (%s)",
	return fmt.Sprintf("%s %.4s [%s] %s",
		fmt.Sprint(rec.Time)[:19],
		alog.LevelNames[rec.Level],
		//
		// XXX TODO
		// u nekom trenu poceo je ispisivati full path do filea sto je kriticno necitko
		// [C:\Users\Neven\Dropbox\Seven\projects\go\workspace\src\nivas.hr\chatprinter\chatprinter.exe]
		// mozda je do gorc1.4
		//		rec.LoggerName,
		//"chatprinter", // hardkodiram za log
		filepath.Base(rec.Filename)+":"+strconv.Itoa(rec.Line),
		rec.Message)

	//		rec.Message+", "+strconv.Itoa(rec.ProcessID)+", "+rec.ProcessName)	// 302100, chatapp.exe
}

// original
// // Format outputs a message like "2014-02-28 18:15:57 [example] INFO     somethinfig happened"
// func (f logFormatter) Format(rec *alog.Record) string {
// 	return fmt.Sprintf("%s %-8s [%s] %-8s %s",
// 		fmt.Sprint(rec.Time)[:19],
// 		alog.LevelNames[rec.Level],
// 		rec.LoggerName,
// 		filepath.Base(rec.Filename)+":"+strconv.Itoa(rec.Line),
// 		rec.Message
// }

func logFile(fileName string) *os.File {

	if file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660); err == nil {
		return file
	} else {
		//return nil
		log.Fatalf("Cannot open log file '%s': %v\n", fileName, err)
		return nil
	}
}

//https://github.com/cenkalti/rain/blob/d45493ebdc299b1d88a8e3ddd6b0195f6bc1c2ac/log.go
// func recoverAndLog(l logger) {
// 	if err := recover(); err != nil {
// 		buf := make([]byte, 10000)
// 		l.Critical(err, "\n", string(buf[:runtime.Stack(buf, false)]))
// 	}
// }
