package config

import (
	"os"

	"github.com/drodil/envssh/util"
)

type EnvVariables struct {
	Static map[string]string `yaml:"static"`
	Moved  []string          `yaml:"moved"`
}

type File struct {
	Local  string `yaml:"local"`
	Remote string `yaml:"remote"`
}

type ServerConfig struct {
	Host  string       `yaml:"host"`
	Port  uint8        `yaml:"port"`
	Env   EnvVariables `yaml:"env"`
	Files []File       `yaml:"files"`
}

type GlobalConfig struct {
	Env   EnvVariables `yaml:"env"`
	Files []File       `yaml:"files"`
}

type Config struct {
	Global  GlobalConfig   `yaml:"global"`
	Servers []ServerConfig `yaml:"servers"`
}

func (config *Config) GetServerConfig(remote *util.Remote) *ServerConfig {
	for _, conf := range config.Servers {
		if conf.Host == remote.Hostname {
			return &conf
		}
	}
	return nil
}

func (config *Config) GetFilesForRemote(remote *util.Remote) []File {
	ret := make([]File, len(config.Global.Files))
	copy(ret, config.Global.Files)

	serverConf := config.GetServerConfig(remote)
	if serverConf != nil {
		ret = append(ret, serverConf.Files...)
	}

	return ret
}

func (config *Config) GetEnvironmentVariablesForRemote(remote *util.Remote) map[string]string {
	ret := make(map[string]string)
	for name, value := range config.Global.Env.Static {
		ret[name] = value
	}

	for _, name := range config.Global.Env.Moved {
		value, ok := os.LookupEnv(name)
		if ok {
			ret[name] = value
		}
	}

	serverConf := config.GetServerConfig(remote)
	if serverConf != nil {
		for name, value := range serverConf.Env.Static {
			ret[name] = value
		}

		for _, name := range serverConf.Env.Moved {
			value, ok := os.LookupEnv(name)
			if ok {
				ret[name] = value
			}
		}
	}

	return ret
}
