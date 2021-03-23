package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/drodil/envssh/util"
	"gopkg.in/yaml.v2"
)

// DefaultConfig is the default configuration for the tool that is written
// to the YAML file on first connect if it does not exist.
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
			Host:    "localhost",
			Port:    22,
			Aliases: []string{"127.0.0.1"},
			Env: EnvVariables{
				Static: map[string]string{},
				Moved:  []string{},
			},
			Files:    []File{},
			Commands: []string{},
		},
	},
}

// ParseConfig parses envssh configuration from given location. Creates new
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

// GetDefaultConfigLocation returns default envssh configuration file
// location.
func GetDefaultConfigLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "envssh.yml"
	}

	return filepath.Join(homeDir, ".ssh", "envssh.yml")
}

// CreateDefaultConfigFile creates default configuration file
// to given location.
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
