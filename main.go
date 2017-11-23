package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/asdine/storm"
	"github.com/namsral/flag"
	"github.com/thoj/go-ircevent"
)

const (
	nickname = "cadmus"
	username = "Cadmus"
	realname = "Founder, King"
)

var (
	db *storm.DB
)

type Addr struct {
	Host   string
	Port   int
	UseTLS bool
}

func (a *Addr) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func ParseAddr(s string) (addr *Addr, err error) {
	addr = &Addr{}

	parts := strings.Split(s, ":")
	fmt.Printf("%v", parts)
	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed address: %s", s)
	}

	addr.Host = parts[0]

	if parts[1][0] == '+' {
		port, err := strconv.Atoi(parts[1][1:])
		if err != nil {
			return nil, fmt.Errorf("invalid port: %s", parts[1])
		}
		addr.Port = port
		addr.UseTLS = true
	} else {
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid port: %s", parts[1])
		}
		addr.Port = port
	}

	if addr.Port < 1 || addr.Port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", addr.Port)
	}

	return addr, nil
}

func main() {
	var (
		err error

		version bool
		config  string
		debug   bool

		dbpath string
	)

	flag.BoolVar(&version, "v", false, "display version information")
	flag.StringVar(&config, "c", "", "config file")
	flag.BoolVar(&debug, "d", false, "debug logging")

	flag.StringVar(&dbpath, "dbpath", "cadmus.db", "path to database")

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

	conn := irc.IRC(nickname, username)
	conn.RealName = realname

	conn.VerboseCallbackHandler = debug
	conn.Debug = debug

	conn.UseTLS = addr.UseTLS
	conn.KeepAlive = 30 * time.Second
	conn.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn.AddCallback("001", func(e *irc.Event) {
		log.Info("Connected!")

		var channels []Channel

		err := db.All(&channels)
		if err != nil {
			log.Fatalf("error loading channels from db: %s", err)
		}

		for _, channel := range channels {
			conn.Join(channel.Name)
			conn.Mode(channel.Name)
			log.Infof("Joined %s", channel.Name)
		}
	})

	conn.AddCallback("INVITE", func(e *irc.Event) {
		var channel Channel
		err = db.One("Name", e.Arguments[0], &channel)
		if err != nil && err == storm.ErrNotFound {
			channel = NewChannel(e.Arguments[1])
			err := db.Save(&channel)
			if err != nil {
				log.Fatalf("error saving channel to db: %s", err)
			}
		} else if err != nil {
			log.Fatalf("error looking up channel in db: %s", err)
		}

		log.Infof("Requested to join %s", channel.Name)

		conn.Join(channel.Name)
		conn.Mode(channel.Name)
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		if e.Arguments[0][0] != '#' {
			return
		}

		log.Infof("[%s] <%s> %s", time.Now().Format("15:04:05"), e.User, e.Message())
	})

	err = conn.Connect(addr.String())
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	conn.Loop()
}
