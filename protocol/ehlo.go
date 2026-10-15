package protocol

import (
	"bufio"
	"regexp"
	"smtp_honeypot/settings"
	"strings"
)

var ehloRegex = regexp.MustCompile("(?i)(EHLO$)|(EHLO (.*))")

type SmtpEhloMessage struct {
}

func (s SmtpEhloMessage) Matches(arg []byte) bool {
	return ehloRegex.Match(arg)
}

func (s SmtpEhloMessage) Handle(connection *SmtpConnection, arg []byte, _ *bufio.Reader) string {
	matches := ehloRegex.FindSubmatch(arg)
	if len(matches) >= 4 {
		connection.Hostname = string(matches[3])
	}

	connection.ExtendedProtocol = true

	identity := settings.GetSettings().HoneypotIdentity.GetValue(connection.SrcAddr)
	return strings.ReplaceAll("250-"+identity+`
250-PIPELINING
250-SIZE 10240000
250-VRFY
250-ETRN
250-STARTTLS
250-AUTH LOGIN PLAIN
250 8BITMIME`, "\n", "\r\n")
}
