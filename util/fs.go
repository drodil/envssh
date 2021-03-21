package util

import (
	"os"
)

// Checks if file exists in given location.
func FileExists(filename string) bool {
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}
