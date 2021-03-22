package util

import (
	"reflect"
	"testing"
)

func TestParseRemote(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want *Remote
	}{
		{name: "Should have host only", str: "localhost", want: &Remote{Hostname: "localhost", Port: uint16(22), Username: ""}},
		{name: "Should have host and port", str: "localhost:1234", want: &Remote{Hostname: "localhost", Port: uint16(1234), Username: ""}},
		{name: "Should have host, port and username", str: "drodil@localhost:1234", want: &Remote{Hostname: "localhost", Port: uint16(1234), Username: "drodil"}},
		{name: "Should have host, port and username with at", str: "drodil@abc@localhost:1234", want: &Remote{Hostname: "localhost", Port: uint16(1234), Username: "drodil@abc"}},
		{name: "Should have default values", str: "", want: &Remote{Hostname: "", Port: uint16(22), Username: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseRemote(tt.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRemote() = %v, want %v", got, tt.want)
			}
		})
	}
}
