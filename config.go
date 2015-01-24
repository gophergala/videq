package main

import (
	"code.google.com/p/gcfg"
	"os"
	"strings"
)

// Config structure that holds Config data
type Config struct {
	FILES struct {
		TEMPLATES string
	}

	DB struct {
		HOST  string
		NAME  string
		USER  string
		PASS  string
		DEBUG bool
	}
	HTTP struct {
		HOSTNAME      string
		LISTENADDRESS string
	}
}

var config Config

// LoadConfig fills Config struct with file data
func LoadConfig() {

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Cannot read hostname: ", err)
	}
	hostname = strings.ToLower(hostname)

	if err := readConfig("./conf/"+hostname+".config.ini", &config); err != nil {
		log.Fatal("Cannot read config file: ", err)
	}
	// log.Info("Loaded config.ini")
	// log.Debug("%#v\n\n", &config)
}
func readConfig(file string, config *Config) error {
	return gcfg.ReadFileInto(config, file)
}
