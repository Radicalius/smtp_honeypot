package protocol

import (
	"fmt"
	"regexp"
)

var heloRegex = regexp.MustCompile("HELO (.*)")

type SmtpHeloMessage struct{}

func (s SmtpHeloMessage) Matches(arg string) bool {
	return heloRegex.MatchString(arg)
}

func (s SmtpHeloMessage) Handle(transaction *SmtpTransaction, arg string) string {
	matches := heloRegex.FindStringSubmatch(arg)
	if len(matches) >= 2 {
		transaction.Hostname = matches[1]
	} else {
		fmt.Println("warning: helo has no hostname submatch")
	}

	return "250 example.com"
}
