package main

import (
	"crypto/tls"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/asdine/storm"
	"github.com/thoj/go-ircevent"
)

type Config struct {
	nick string
	user string
	name string

	addr   *Addr
	debug  bool
	logdir string
}

type Bot struct {
	config *Config
	conn   *irc.Connection

	network string
	loggers map[string]Logger
}

func NewBot(config *Config) *Bot {
	return &Bot{
		config:  config,
		loggers: make(map[string]Logger),
	}
}

func (b *Bot) getLogger(channel string) (logger Logger, err error) {
	logger = b.loggers[channel]
	if logger == nil {
		logger, err = NewFileLogger(b.config.logdir, b.network, channel)
		if err != nil {
			return nil, err
		}
		b.loggers[channel] = logger
	}
	return
}

func (b *Bot) onInvite(e *irc.Event) {
	var channel Channel
	err := db.One("Name", e.Arguments[0], &channel)
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

	b.conn.Join(channel.Name)
	b.conn.Mode(channel.Name)
}

func (b *Bot) onConnected(e *irc.Event) {
	var channels []Channel

	log.Info("Connected!")

	p := regexp.MustCompile("^[Ww]elcome to the (.*) Internet Relay Network")
	matches := p.FindStringSubmatch(e.Message())
	if len(matches) == 2 {
		b.network = matches[1]
		log.Infof("Network is %s", b.network)
	}

	err := db.All(&channels)
	if err != nil {
		log.Fatalf("error loading channels from db: %s", err)
	}

	for _, channel := range channels {
		b.conn.Join(channel.Name)
		b.conn.Mode(channel.Name)
		log.Infof("Joined %s", channel.Name)
	}
}

func (b *Bot) onMessage(e *irc.Event) {
	if e.Arguments[0][0] != '#' {
		return
	}

	channel := e.Arguments[0]

	logger, err := b.getLogger(channel)
	if err != nil {
		log.Errorf(
			"error getting logger for %s on %s: %s",
			channel, b.network, err,
		)
		return
	}

	logger.LogMessage(e.User, e.Message())
}

func (b *Bot) Run() error {
	b.conn = irc.IRC(b.config.nick, b.config.user)
	b.conn.RealName = b.config.name

	b.conn.VerboseCallbackHandler = b.config.debug
	b.conn.Debug = b.config.debug

	b.setupCallbacks()

	b.conn.UseTLS = b.config.addr.UseTLS
	b.conn.KeepAlive = 30 * time.Second
	b.conn.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := b.conn.Connect(b.config.addr.String())
	if err != nil {
		return err
	}

	b.conn.Loop()
	return nil
}

func (b *Bot) setupCallbacks() {
	b.conn.AddCallback("001", b.onConnected)
	b.conn.AddCallback("INVITE", b.onInvite)
	b.conn.AddCallback("PRIVMSG", b.onMessage)
}
