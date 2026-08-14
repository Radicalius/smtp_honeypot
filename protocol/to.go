package protocol

import (
	"fmt"
	"regexp"
)

var toRegex = regexp.MustCompile("RCPT TO:<([^>]*)>")

type SmtpToMessage struct {
}

func (s SmtpToMessage) Matches(arg string) bool {
	return toRegex.MatchString(arg)
}

func (s SmtpToMessage) Handle(transaction *SmtpTransaction, arg string) string {
	matches := toRegex.FindStringSubmatch(arg)
	if len(matches) < 2 {
		fmt.Println("warning: no to match")
	} else if len(transaction.To) >= 100 {
		fmt.Println("warning: more than 100 to fields supplied")
	} else {
		transaction.To = append(transaction.To, matches[1])
	}

	return "250 2.1.5 OK"
}
