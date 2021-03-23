package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/drodil/envssh/config"
	"github.com/drodil/envssh/ssh"
	"github.com/drodil/envssh/util"
)

var logger = util.GetLogger()

func main() {

	// TODO: Move handling of params to own function
	username := flag.String("l", util.GetUsername(), "Login username")
	configFile := flag.String("c", config.GetDefaultConfigLocation(), "envssh configuration file")
	port := flag.Uint("p", 22, "Port to connect to")

	flag.Parse()

	destination := flag.Arg(0)
	if destination == "" {
		// TODO: Explain usage better (destination, port, etc.)
		flag.Usage()
		os.Exit(0)
	}

	cmd := ""
	nArgs := flag.NArg()
	if nArgs > 1 {
		for i := 1; i < nArgs; i++ {
			cmd = fmt.Sprint(cmd, " ", flag.Arg(i))
		}
	}

	remote := util.ParseRemote(destination)
	if remote.Username == "" {
		remote.Username = *username
	}

	if *port != uint(22) {
		remote.Port = uint16(*port)
	}

	config, err := config.ParseConfig(*configFile)
	if err != nil {
		logger.Fatal(err)
		fmt.Println("Failed to parse configuration file")
		os.Exit(1)
	}

	serverConf := config.GetServerConfig(remote)
	if serverConf != nil {
		if serverConf.Port != 0 {
			remote.Port = serverConf.Port
		}
		remote.Hostname = serverConf.Host
	}

	client, err := ssh.ConnectAuto(remote)
	if err != nil {
		logger.Fatal(err)
		fmt.Println("Disconnected from", remote.Hostname, "port", remote.Port)
		os.Exit(1)
	}

	err = setUpRemote(client, config, remote)
	if err != nil {
		logger.Fatal(err)
		fmt.Println("Connection failed to ", remote.ToAddress(), err)
		os.Exit(1)
	}

	if cmd == "" {
		err = client.StartInteractiveSession()
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		err = client.RunCommand(cmd, nil, os.Stdout, os.Stderr)
		if err != nil {
			logger.Fatal(err)
		}
	}
	client.Disconnect()
}

func setUpRemote(client *ssh.Client, config *config.Config, remote *util.Remote) error {
	client.SetRemoteEnvMap(config.GetEnvironmentVariablesForRemote(remote))

	for _, file := range config.GetFilesForRemote(remote) {
		local := os.ExpandEnv(file.Local)
		// TODO: Make this configurable?
		if !util.FileExists(local) {
			logger.Println("Local file", local, "missing, skipping copy to remote")
			continue
		}

		err := client.CopyFileToRemote(local, file.Remote)
		if err != nil {
			return err
		}
	}

	for _, cmd := range config.GetCommandsForRemote(remote) {
		err := client.RunCommand(cmd, nil, nil, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
