package config

import (
	"reflect"
	"testing"

	"github.com/drodil/envssh/util"
)

func TestConfig_GetServerConfig(t *testing.T) {
	serverConfig := &ServerConfig{
		Host:    "126.0.0.1",
		Aliases: []string{"125.0.0.1"},
	}

	otherServerConfig := &ServerConfig{
		Host: "10.40.124",
	}

	config := &Config{
		Servers: []ServerConfig{
			*serverConfig,
			*otherServerConfig,
		},
	}

	tests := []struct {
		name   string
		config *Config
		remote *util.Remote
		want   *ServerConfig
	}{
		{
			name:   "Should return by hostname",
			config: config,
			remote: &util.Remote{
				Hostname: "126.0.0.1",
			},
			want: serverConfig,
		},
		{
			name:   "Should return by alias",
			config: config,
			remote: &util.Remote{
				Hostname: "125.0.0.1",
			},
			want: serverConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.GetServerConfig(tt.remote); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.GetServerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
