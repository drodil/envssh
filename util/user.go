package util

import (
	"os/user"
)

func GetUsername() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Username
}
