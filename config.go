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

/*
https://github.com/Xe/cqbot/blob/70a909d601f144e997a72867e48e89ca93a8668f/bot/bot.go

type Bot struct {
	server   string
	port     string
	nick     string
	user     string
	gecos    string
	handlers map[string]map[string]*Handler
	Commands map[string]*Command
	Channels map[string]*Channel
	Scripts  map[string]*Script
	Config   Config
	Conn     net.Conn
	Log      *log.Logger
}

// Allocate the maps and seed the configuration details of the bot
// from the configuration file. Also specify handlers for common IRC
// protocol verbs.
func NewBot(confname string) (bot *Bot) {
	bot = &Bot{
		handlers: make(map[string]map[string]*Handler),
		Config:   LoadConfig(confname),
		Commands: make(map[string]*Command),
		Channels: make(map[string]*Channel),
		Scripts:  make(map[string]*Script),
		Log:      log.New(os.Stdout, "", log.LstdFlags),
	}
}



https://github.com/Xe/cqbot/blob/70a909d601f144e997a72867e48e89ca93a8668f/bot/config.go

type ServerConfig struct {
	Port string
	Host string
}

type BotConfig struct {
	Nick    string
	User    string
	Gecos   string
	Channel string
	Nspass  string
	Prefix  string
}

type Config struct {
	Server ServerConfig
	Bot    BotConfig
}

// Wrap the ini file and load the configuration
func LoadConfig(cfgFile string) (cfg Config) {
	err := gcfg.ReadFileInto(&cfg, cfgFile)

	if err != nil {
		panic(err)
	}

	return
}

func writeLine(fout *os.File, line string) {
	num, err := fout.Write([]byte(line + "\n"))
	if num != len(line) {
		return
	}
	if err != nil {
		panic(err)
	}
}

func (conf *Config) Export(fname string) (err error) {
	fout, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	writeLine(fout, "[server]")
	writeLine(fout, "port = " + conf.Server.Port)
	writeLine(fout, "host = " + conf.Server.Host)
	writeLine(fout, "\n[bot]")
	writeLine(fout, "nick = " + conf.Bot.Nick)
	writeLine(fout, "user = " + conf.Bot.User)
	writeLine(fout, "gecos = " + conf.Bot.Gecos)
	writeLine(fout, "channel = " + conf.Bot.Channel)
	writeLine(fout, "nspass = " + conf.Bot.Nspass)
	writeLine(fout, "prefix = " + conf.Bot.Prefix)

	return
}
*/
