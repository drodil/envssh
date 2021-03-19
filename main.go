package main

import (
	"flag"
	"github.com/drodil/envssh/ssh"
	"os/user"
)

func main() {

	username := flag.String("l", getUsername(), "Login username")

	flag.Parse()

	destination := flag.Arg(0)

	if destination == "" {
		flag.Usage()
		return
	}

	client, err := ssh.ConnectWithPassword(destination, *username)
	if err != nil {
		panic(err)
	}

	client.Disconnect()
}

// TODO: Move to utils
func getUsername() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Username
}
