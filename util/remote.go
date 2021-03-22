package util

import (
	"fmt"
	"strconv"
	"strings"
)

// Remote presents a SSH remote with username, hostname and port.
type Remote struct {
	Username string
	Hostname string
	Port     uint8
}

// ParseRemote parses SSH remote from string. Takes into account username before last
// @ character and port after : character.
func ParseRemote(str string) *Remote {
	var port uint8
	port = 22
	hostname := str
	username := ""
	if strings.Contains(str, "@") {
		hostname = str[strings.LastIndex(str, "@")+1:]
		username = str[:strings.LastIndex(str, "@")]
		str = hostname
	}

	if strings.Contains(str, ":") {
		parts := strings.Split(str, ":")
		hostname = parts[0]
		portConv, err := strconv.ParseUint(parts[1], 10, 8)
		if err == nil {
			port = uint8(portConv)
		}
	}
	return &Remote{Username: username, Hostname: hostname, Port: port}
}

// ToAddress returns hostname:port presentation of Remote.
func (remote *Remote) ToAddress() string {
	return fmt.Sprint(remote.Hostname, ":", strconv.Itoa(int(remote.Port)))
}
