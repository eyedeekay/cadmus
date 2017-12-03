package cadmus

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Logger interface {
	Rotate() error
	Log(message string) error
	Logf(format string, args ...interface{}) error

	Channel() string
	Network() string

	LogMessage(user, message string) error
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
	err := os.MkdirAll(pathname, 0755)
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

	log.Infof("Logger created for %s on %s", channel, network)

	return &FileLogger{
		logdir:  logdir,
		network: network,
		channel: channel,
		f:       f,
	}, nil
}

func (l *FileLogger) Channel() string {
	return l.channel
}

func (l *FileLogger) Network() string {
	return l.network
}

func (l *FileLogger) Rotate() error {
	l.Lock()
	defer l.Unlock()

	l.f.Close()

	filename := path.Join(
		l.logdir, l.network, l.channel,
		fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")),
	)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l.f = f

	return nil
}

func (l *FileLogger) Log(message string) error {
	l.Lock()
	defer l.Unlock()

	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	_, err := l.f.WriteString(message)
	return err
}

func (l *FileLogger) Logf(format string, args ...interface{}) error {
	return l.Log(fmt.Sprintf(format, args...))
}

func (l *FileLogger) LogMessage(user, message string) error {
	ts := time.Now().Format("15:04:05")
	return l.Logf("[%s] <%s> %s", ts, user, message)
}
