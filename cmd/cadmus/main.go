package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/namsral/flag"

	"github.com/prologic/cadmus"
)

const (
	nick = "cadmus"
	user = "Cadmus"
	name = "Founder, King"
)

func main() {
	var (
		version bool
		config  string
		debug   bool

		dbpath  string
		logpath string
	)

	flag.BoolVar(&version, "v", false, "display version information")
	flag.StringVar(&config, "c", "", "config file")
	flag.BoolVar(&debug, "d", false, "debug logging")

	flag.StringVar(&dbpath, "dbpath", "cadmus.db", "path to database")
	flag.StringVar(&logpath, "logpath", "./logs", "path to store logs")

	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if version {
		fmt.Printf(cadmus.FullVersion())
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		log.Fatalf("Ussage: %s <address>[:port]", os.Args[0])
	}

	bot := cadmus.NewBot(flag.Arg(0), &cadmus.Config{
		Nick:    nick,
		User:    user,
		Name:    name,
		Debug:   debug,
		DBPath:  dbpath,
		LogPath: logpath,
	})
	log.Fatal(bot.Run())
}
