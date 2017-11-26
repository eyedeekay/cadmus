package main

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

type Logger interface {
	Rotate() error
	Log(message string) error
}

type FileLogger struct {
	sync.Mutex

	logdir  string
	network string
	channel string

	f *os.File
}

func NewFileLogger(logdir, network, channel string) (*FileLogger, error) {
	pathname := path.Join(logdir, network, channel)
	err := os.MkdirAll(pathname, 0x755)
	if err != nil {
		return nil, err
	}

	filename := path.Join(
		pathname, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")),
	)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FileLogger{
		logdir:  logdir,
		network: network,
		channel: channel,
		f:       f,
	}, nil
}

func (l *FileLogger) Rorate() error {
	l.Lock()
	defer l.Unlock()

	l.f.Close()

	logfile := path.Join(
		l.logdir,
		fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")),
	)

	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l.f = f

	return nil
}

func (l *FileLogger) Log(message string) error {
	l.Lock()
	defer l.Unlock()

	_, err := l.f.WriteString(message)
	return err
}
