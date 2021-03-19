package util

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

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
	    answer = strings.Replace(answer, "\n", "", -1)
        if len(answer) == 0 {
            answer = def
            break
        }

        if !StringInSlice(answer, allowed) {
            allowedStr := strings.Join(allowed, "', ")
            question = fmt.Sprint("Please type '", allowedStr, "':")
	    }
	    answer = strings.Replace(answer, "\n", "", -1)
        break
    }

    return answer
}
