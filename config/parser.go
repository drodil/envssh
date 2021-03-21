package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/drodil/envssh/util"
	"gopkg.in/yaml.v2"
)

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

func GetDefaultConfigLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "envssh.yml"
	}

	return filepath.Join(homeDir, ".ssh", "envssh.yml")
}

func CreateDefaultConfigFile(location string) error {
	def := &Config{
		Global: GlobalConfig{
			Env: EnvVariables{
				Static: map[string]string{"LC_ENVSSH": "1"},
				Moved:  []string{"LANG", "EDITOR", "VISUAL"},
			},
			Files:    []File{{Local: ".bashrc", Remote: ".bashrc"}},
			Commands: []string{"export ENVSSH=1"},
		},
		Servers: []ServerConfig{{Host: "localhost", Port: 22}},
	}

	d, err := yaml.Marshal(&def)
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
