package protocol

import (
	"regexp"
	"strings"
)

var ehloRegex = regexp.MustCompile("(EHLO$)|(EHLO (.*))")

type SmtpEhloMessage struct {
}

func (s SmtpEhloMessage) Matches(arg []byte) bool {
	return ehloRegex.Match(arg)
}

func (s SmtpEhloMessage) Handle(connection *SmtpConnection, arg []byte) string {
	matches := ehloRegex.FindSubmatch(arg)
	if len(matches) >= 4 {
		connection.Hostname = string(matches[3])
	}

	connection.ExtendedProtocol = true
	return strings.ReplaceAll(`250-{server_name}
250-PIPELINING
250-SIZE 10240000
250-VRFY
250-ETRN
250-STARTTLS
250-AUTH LOGIN PLAIN
250 8BITMIME`, "\n", "\r\n")
}
