package protocol

import (
	"fmt"
	"regexp"
)

var fromRegex = regexp.MustCompile("MAIL FROM:<([^>]*)>")

type SmtpFromMessage struct {
}

func (s SmtpFromMessage) Matches(arg string) bool {
	return fromRegex.MatchString(arg)
}

func (s SmtpFromMessage) Handle(transaction *SmtpTransaction, arg string) string {
	matches := fromRegex.FindStringSubmatch(arg)
	if len(matches) < 2 {
		fmt.Println("warning: no from match")
	} else if len(transaction.From) >= 100 {
		fmt.Println("warning: more than 100 from fields supplied")
	} else {
		transaction.From = append(transaction.From, matches[1])
	}

	return "250 2.1.0 OK"
}
