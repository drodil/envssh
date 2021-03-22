package config

import (
	"os"

	"github.com/drodil/envssh/util"
)

// EnvVariables that should be set on remote.
type EnvVariables struct {
	Static map[string]string `yaml:"static"`
	Moved  []string          `yaml:"moved"`
}

// File presents a single file that is moved to remote.
type File struct {
	Local  string `yaml:"local"`
	Remote string `yaml:"remote"`
}

// ServerConfig presents a single server configuration
// identified by the host name.
type ServerConfig struct {
	Host     string       `yaml:"host"`
	Port     uint16        `yaml:"port"`
	Env      EnvVariables `yaml:"env"`
	Files    []File       `yaml:"files"`
	Commands []string     `yaml:"commands"`
}

// GlobalConfig is base for all server configurations
// and is used for all remotes.
type GlobalConfig struct {
	Env      EnvVariables `yaml:"env"`
	Files    []File       `yaml:"files"`
	Commands []string     `yaml:"commands"`
}

// Config contains global and server specific configurations.
type Config struct {
	Global  GlobalConfig   `yaml:"global"`
	Servers []ServerConfig `yaml:"servers"`
}

// GetServerConfig returns server specific config based on hostname if it exists in the Config struct.
func (config *Config) GetServerConfig(remote *util.Remote) *ServerConfig {
	for _, conf := range config.Servers {
		if conf.Host == remote.Hostname {
			return &conf
		}
	}
	return nil
}

// GetCommandsForRemote returns a list of commands to run on remote.
func (config *Config) GetCommandsForRemote(remote *util.Remote) []string {
	ret := make([]string, len(config.Global.Commands))
	copy(ret, config.Global.Commands)
	serverConf := config.GetServerConfig(remote)
	if serverConf != nil {
		ret = append(ret, serverConf.Commands...)
	}
	return ret
}

// GetFilesForRemote returns a list of files that are transferred to remote.
func (config *Config) GetFilesForRemote(remote *util.Remote) []File {
	ret := make([]File, len(config.Global.Files))
	copy(ret, config.Global.Files)

	serverConf := config.GetServerConfig(remote)
	if serverConf != nil {
		ret = append(ret, serverConf.Files...)
	}

	return ret
}

// GetEnvironmentVariablesForRemote returns a list of environment variables that are set to remote.
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
