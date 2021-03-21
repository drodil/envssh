package main

import (
	"flag"

	"github.com/drodil/envssh/config"
	"github.com/drodil/envssh/ssh"
	"github.com/drodil/envssh/util"
)

var logger = util.GetLogger()

func main() {

	// TODO: Move handling of params to own function
	username := flag.String("l", util.GetUsername(), "Login username")
	configFile := flag.String("c", config.GetDefaultConfigLocation(), "envssh configuration file")

	flag.Parse()

	destination := flag.Arg(0)
	remote := util.ParseRemote(destination)
	if remote.Username == "" {
		remote.Username = *username
	}

	if destination == "" {
		// TODO: Explain usage better (destination, port, etc.)
		flag.Usage()
		return
	}

	config, err := config.ParseConfig(*configFile)
	if err != nil {
		panic(err)
	}

	serverConf := config.GetServerConfig(remote)
	if serverConf != nil && serverConf.Port != 0 {
		remote.Port = serverConf.Port
	}

	client, err := ssh.ConnectWithPassword(remote)
	if err != nil {
		panic(err)
	}

	err = setUpRemote(client, config, remote)
	if err != nil {
		panic(err)
	}

	err = client.StartInteractiveSession()
	if err != nil {
		panic(err)
	}
	client.Disconnect()
}

func setUpRemote(client *ssh.Client, config *config.Config, remote *util.Remote) error {
	client.SetRemoteEnvMap(config.GetEnvironmentVariablesForRemote(remote))

	for _, file := range config.GetFilesForRemote(remote) {
		// TODO: Make this configurable?
		if !util.FileExists(file.Local) {
			logger.Println("Local file", file.Local, "missing, skipping copy to remote")
			continue
		}

		err := client.CopyFileToRemote(file.Local, file.Remote)
		if err != nil {
			return err
		}
	}

	for _, cmd := range config.GetCommandsForRemote(remote) {
		err := client.RunCommand(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
