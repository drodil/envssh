package main

import (
	"flag"
	"os/user"
	"github.com/drodil/envssh/ssh"
)

func main() {

    username := flag.String("l", getUsername(), "Login username")

    flag.Parse()

    destination := flag.Arg(0)

    if destination == "" {
        flag.Usage()
        return
    }

    client, err := ssh.ConnectWithPassword(destination, *username, "password")
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
