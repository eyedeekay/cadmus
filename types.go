package cadmus

import (
	"fmt"
	"strings"
	"sync"
)

// ChannelLoggerMap holds a mapping of channels to Logger interfaces
// that is safe for concurrent readers and writers.
type ChannelLoggerMap struct {
	sync.RWMutex
	channels map[string]Logger
}

// NewChannelLoggerMap returns a new initialized *ChannelLoggerMap
func NewChannelLoggerMap() *ChannelLoggerMap {
	return &ChannelLoggerMap{
		channels: make(map[string]Logger),
	}
}

// Count returns the number of Logger(s)
func (c *ChannelLoggerMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.channels)
}

// Range ranges over the Logger(s) calling f
func (c *ChannelLoggerMap) Range(f func(kay string, value Logger) bool) {
	c.Lock()
	defer c.Unlock()
	for k, v := range c.channels {
		if !f(k, v) {
			return
		}
	}
}

// Get returns a Logger given a channel if it exists or a zero-value Logger
func (c *ChannelLoggerMap) Get(channel string) Logger {
	c.RLock()
	defer c.RUnlock()
	return c.channels[strings.ToLower(channel)]
}

// Add adds a new *Channel if not already exists or an error otherwise
func (c *ChannelLoggerMap) Add(logger Logger) error {
	c.Lock()
	defer c.Unlock()
	if c.channels[strings.ToLower(logger.Channel())] != nil {
		return fmt.Errorf("%s already exists", logger.Channel())
	}
	c.channels[strings.ToLower(logger.Channel())] = logger
	return nil
}
