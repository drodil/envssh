package main

import (
	"flag"
	"os"

	"github.com/drodil/envssh/config"
	"github.com/drodil/envssh/ssh"
	"github.com/drodil/envssh/util"
)

func main() {

	username := flag.String("l", util.GetUsername(), "Login username")
	configFile := flag.String("c", config.GetDefaultConfigLocation(), "envssh configuration file")

	flag.Parse()

	destination := flag.Arg(0)
	remote := ssh.ParseRemote(destination)
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

	setEnvFromConfig(client, config)
	err = client.StartInteractiveSession()
	if err != nil {
		panic(err)
	}
	client.Disconnect()
}

func setEnvFromConfig(client *ssh.Client, config *config.Config) {
	for name, value := range config.Global.Env.Static {
		client.SetRemoteEnv(name, value)
	}

	for _, name := range config.Global.Env.Moved {
		value, ok := os.LookupEnv(name)
		if ok {
			client.SetRemoteEnv(name, value)
		}
	}
}
