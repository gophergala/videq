package main

import (
	"os/exec"
)

var execList = map[string]string{
	"mediainfo": "sudo apt-get install mediainfo",
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
