package util

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	file, err := ioutil.TempFile("", "testFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{name: "Should exist", filename: file.Name(), want: true},
		{name: "Should not exist", filename: "invalid", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.filename); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
