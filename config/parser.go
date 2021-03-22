package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/drodil/envssh/util"
	"gopkg.in/yaml.v2"
)

var DefaultConfig = &Config{
	Global: GlobalConfig{
		Env: EnvVariables{
			Static: map[string]string{"LC_ENVSSH": "1"},
			Moved:  []string{"LANG", "EDITOR", "VISUAL"},
		},
		Files:    []File{{Local: "$HOME/.bashrc", Remote: "$HOME/.bashrc"}},
		Commands: []string{"export ENVSSH=1"},
	},
	Servers: []ServerConfig{
		{
			Host: "localhost",
			Port: 22,
			Env: EnvVariables{
				Static: map[string]string{},
				Moved:  []string{},
			},
			Files:    []File{},
			Commands: []string{},
		},
	},
}

// Parses sshenv configuration from given location. Creates new
// default configuration to the location if it does not exist.
func ParseConfig(location string) (*Config, error) {
	if !util.FileExists(location) {
		fmt.Println("enssh configuration missing, creating default to", location)
		err := CreateDefaultConfigFile(location)
		if err != nil {
			return nil, err
		}
	}

	confFile, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(confFile, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Returns default envssh configuration location.
func GetDefaultConfigLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "envssh.yml"
	}

	return filepath.Join(homeDir, ".ssh", "envssh.yml")
}

// Creates default configuration to given location.
func CreateDefaultConfigFile(location string) error {
	d, err := yaml.Marshal(&DefaultConfig)
	if err != nil {
		return err
	}

	f, err := os.Create(location)
	if err != nil {
		return nil
	}
	defer f.Close()
	return ioutil.WriteFile(location, d, 0744)
}
