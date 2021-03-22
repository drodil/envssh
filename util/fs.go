package util

import (
	"os"
)

// FileExists checks if file exists in given location.
func FileExists(filename string) bool {
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}
