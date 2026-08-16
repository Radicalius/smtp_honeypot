package protocol

import (
	"fmt"
	"regexp"
)

var heloRegex = regexp.MustCompile("HELO (.*)")

type SmtpHeloMessage struct{}

func (s SmtpHeloMessage) Matches(arg []byte) bool {
	return heloRegex.Match(arg)
}

func (s SmtpHeloMessage) Handle(transaction *SmtpTransaction, arg []byte) string {
	matches := heloRegex.FindSubmatch(arg)
	if len(matches) >= 2 {
		transaction.Hostname = string(matches[1])
	} else {
		fmt.Println("warning: helo has no hostname submatch")
	}

	return "250 example.com"
}
