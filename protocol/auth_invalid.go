package protocol

import (
	"regexp"
)

var authInvalid = regexp.MustCompile(`(?i)AUTH`)

type SmtpAuthInvalidMessage struct {
}

func (s SmtpAuthInvalidMessage) Matches(arg []byte) bool {
	return authInvalid.Match(arg)
}

func (s SmtpAuthInvalidMessage) Handle(connection *SmtpConnection, arg []byte) string {
	return "504 5.5.4 Unrecognized authentication type"
}
