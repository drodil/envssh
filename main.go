package main

import (
	"flag"

	"github.com/drodil/envssh/config"
	"github.com/drodil/envssh/ssh"
	"github.com/drodil/envssh/util"
)

func main() {

	username := flag.String("l", util.GetUsername(), "Login username")
	configFile := flag.String("c", config.GetDefaultConfigLocation(), "envssh configuration file")

	flag.Parse()

	destination := flag.Arg(0)
	remote := util.ParseRemote(destination)
	if remote.Username == "" {
		remote.Username = *username
	}

	if destination == "" {
		flag.Usage()
		return
	}

	config, err := config.ParseConfig(*configFile)
	if err != nil {
		panic(err)
	}

	client, err := ssh.ConnectWithPassword(remote)
	if err != nil {
		panic(err)
	}

	client.SetRemoteEnvMap(config.GetEnvironmentVariablesForRemote(remote))
	err = client.StartInteractiveSession()
	if err != nil {
		panic(err)
	}
	client.Disconnect()
}
