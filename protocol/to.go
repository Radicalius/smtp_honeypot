package protocol

import (
	"fmt"
	"regexp"
)

var toRegex = regexp.MustCompile("RCPT TO:<([^>]*)>")

type SmtpToMessage struct {
}

func (s SmtpToMessage) Matches(arg []byte) bool {
	return toRegex.Match(arg)
}

func (s SmtpToMessage) Handle(transaction *SmtpTransaction, arg []byte) string {
	matches := toRegex.FindSubmatch(arg)
	if len(matches) < 2 {
		fmt.Println("warning: no to match")
	} else if len(transaction.To) >= 100 {
		fmt.Println("warning: more than 100 to fields supplied")
	} else {
		transaction.To = append(transaction.To, string(matches[1]))
	}

	return "250 2.1.5 OK"
}
