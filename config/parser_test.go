package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

var data = `
global:
  env:
    static:
      LC_ENVSSH: "1"
    moved:
    - LANG
    - EDITOR
    - VISUAL
  files:
  - local: .bashrc
    remote: .bashrc
  commands:
  - "export TEST=1"
servers:
- host: localhost
  port: 22
  env:
    static:
      LC_ENVSSH: "2"
    moved:
    - LANGUAGE
`

func TestParsing(t *testing.T) {
	file, err := ioutil.TempFile("", "test_config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(data)
	if err != nil {
		t.Fatal(err)
	}

	config, err := ParseConfig(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(config.Global.Env.Static), "Global static env variable count should be one")
	assert.Equal(t, 3, len(config.Global.Env.Moved), "Global moved env variable count should be three")
	assert.Equal(t, 1, len(config.Global.Files), "Global moved files count should be one")
	assert.Equal(t, 1, len(config.Global.Commands), "Global command count should be one")

	assert.Equal(t, 1, len(config.Servers), "Server specific config count should be one")
	assert.Equal(t, "localhost", config.Servers[0].Host, "Hostname of server config should be localhost")
	assert.Equal(t, uint8(22), config.Servers[0].Port, "Port of server config should be 22")
	assert.Equal(t, 1, len(config.Servers[0].Env.Static), "Server specific static env variable count should be one")
	assert.Equal(t, 1, len(config.Servers[0].Env.Moved), "Server specific moved env variable count should be one")
}

func TestDefaultConfig(t *testing.T) {
	file, err := ioutil.TempFile("", "test_config")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(file.Name())

	config, err := ParseConfig(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, DefaultConfig, config, "Default config should be returned")
	os.Remove(file.Name())
}

func TestDefaultConfigLocation(t *testing.T) {
	loc := GetDefaultConfigLocation()
	assert.Contains(t, loc, "envssh.yml")
}
