package main

import (
	"os/exec"
)

var execList = map[string]string{
	"mediainfo":     "\r\nsudo apt-get install mediainfo",
	"HandBrakeCLI":  "\r\nsudo add-apt-repository ppa:stebbins/handbrake-releases\r\napt-get install handbrake-cli\r\n",
	"ffmpeg":        "\r\nsudo add-apt-repository ppa:jon-severinsson/ffmpeg\r\nsudo apt-get update\r\nsudo apt-get install ffmpeg\r\nsudo apt-get install frei0r-plugins\r\n",
	"ffmpeg2theora": "\r\nsudo apt-get install ffmpeg2theora",
}

func checkExecutables() {
	var failed bool = false
	for fileName, installHelp := range execList {
		err := checkExecutable(fileName)
		if err != nil {
			log.Errorln(err)
			log.Info("run: ", installHelp)
			failed = true
		}
	}
	if failed == true {
		log.Fatalln("Missing executables, cannot continue. please fix and rerun.")
	}
}

func checkExecutable(exename string) error {
	_, err := exec.LookPath(exename)
	if err != nil {
		return err
	}
	return nil
}
