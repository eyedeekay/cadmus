package cadmus

import (
	"time"
)

type Channel struct {
	ID        int       `storm:"id,increment"`
	Name      string    `storm:"index"`
	CreatedAt time.Time `storm:"index"`
}

func NewChannel(name string) Channel {
	return Channel{
		Name:      name,
		CreatedAt: time.Now(),
	}
}
