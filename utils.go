package cadmus

import (
	"fmt"
	"strconv"
	"strings"
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
