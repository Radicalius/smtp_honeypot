package protocol

import (
	"fmt"
	"regexp"
)

var ehloRegex = regexp.MustCompile("EHLO (.*)")

type SmtpEhloMessage struct {
}

func (s SmtpEhloMessage) Matches(arg []byte) bool {
	return ehloRegex.Match(arg)
}

func (s SmtpEhloMessage) Handle(transaction *SmtpTransaction, arg []byte) string {
	matches := ehloRegex.FindSubmatch(arg)
	if len(matches) >= 2 {
		transaction.Hostname = string(matches[1])
	} else {
		fmt.Println("warning: ehlo missing submatch")
	}

	transaction.ExtendedProtocol = true
	return `250-{server_name}
250-PIPELINING
250-SIZE 10240000
250-VRFY
250-ETRN
250-STARTTLS
250-AUTH LOGIN PLAIN
250 8BITMIME`
}
