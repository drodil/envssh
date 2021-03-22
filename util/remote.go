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
	Port     uint16
}

// ParseRemote parses SSH remote from string. Takes into account username before last
// @ character and port after : character.
func ParseRemote(str string) *Remote {
	var port uint16
	port = 22
	hostname := str
	username := ""
	if strings.Contains(str, "@") {
		hostname = str[strings.LastIndex(str, "@")+1:]
		username = str[:strings.LastIndex(str, "@")]
		str = hostname
	}

	fmt.Print(str)
	if strings.Contains(str, ":") {
		parts := strings.Split(str, ":")
		hostname = parts[0]
		portConv, err := strconv.ParseUint(parts[1], 10, 16)
		if err == nil {
			port = uint16(portConv)
		}
	}
	return &Remote{Username: username, Hostname: hostname, Port: port}
}

// ToAddress returns hostname:port presentation of Remote.
func (remote *Remote) ToAddress() string {
	return fmt.Sprint(remote.Hostname, ":", strconv.Itoa(int(remote.Port)))
}
