package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// PromptAllowedString prompts user an answer that is in the allowed slice.
// Returns default in case user gives empty answer.
func PromptAllowedString(question string, allowed []string, def string) string {
	var answer string
	var err error
	fmt.Print(question, " ")
	buf := bufio.NewReader(os.Stdin)

	for {
		answer, err = buf.ReadString('\n')
		if err != nil {
			continue
		}
		answer = strings.Replace(strings.Replace(answer, "\r", "", -1), "\n", "", -1)
		if len(answer) == 0 {
			answer = def
			break
		}

		if !stringInSlice(answer, allowed) {
			allowedStr := strings.Join(allowed, "', ")
			question = fmt.Sprint("Please type '", allowedStr, "':")
		}
		answer = strings.Replace(strings.Replace(answer, "\r", "", -1), "\n", "", -1)
		break
	}

	return answer
}

// PromptPassword prompts password from user with given question.
// Password is not shown in terminal.
func PromptPassword(question string) string {
	fmt.Print(question, " ")
	stdIn := int(syscall.Stdin)
	state, _ := terminal.GetState(stdIn)
	defer terminal.Restore(stdIn, state)

	bytePassword, err := terminal.ReadPassword(stdIn)
	fmt.Print("\n")
	if err != nil {
		return ""
	}

	password := string(bytePassword)
	return strings.TrimSpace(password)
}
