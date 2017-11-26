package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/asdine/storm"
	"github.com/namsral/flag"
)

const (
	nick = "cadmus"
	user = "Cadmus"
	name = "Founder, King"
)

var (
	db *storm.DB
)

func main() {
	var (
		err error

		version bool
		config  string
		debug   bool

		dbpath string
		logdir string
	)

	flag.BoolVar(&version, "v", false, "display version information")
	flag.StringVar(&config, "c", "", "config file")
	flag.BoolVar(&debug, "d", false, "debug logging")

	flag.StringVar(&dbpath, "dbpath", "cadmus.db", "path to database")
	flag.StringVar(&logdir, "logdir", "./logs", "path to store logs")

	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if version {
		fmt.Printf("cadmus v%s", FullVersion())
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		log.Fatalf("Ussage: %s <address>[:port]", os.Args[0])
	}

	addr, err := ParseAddr(flag.Arg(0))
	if err != nil {
		log.Fatalf("error parsing addr: %s", err)
	}

	db, err = storm.Open(dbpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bot := NewBot(&Config{
		nick:   nick,
		user:   user,
		name:   name,
		addr:   addr,
		debug:  debug,
		logdir: logdir,
	})
	log.Fatal(bot.Run())
}
