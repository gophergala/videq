package main

import (
	"log"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := createStorage()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

}

func createStorage() error {
	paths := []string{
		"./storage/datastore",
		"./storage/.upload"}

	for _, path := range paths {
		err := os.MkdirAll(path, 02750)
		if err != nil {
			return err
		}
	}

	return nil
}
