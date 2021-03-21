package util

import (
	"os/user"
)

// Returns current username or empty if it cannot be
// fetched from OS.
func GetUsername() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Username
}
