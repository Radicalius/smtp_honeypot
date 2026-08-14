package protocol

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"regexp"
)

var authPlainRegex = regexp.MustCompile("AUTH PLAIN (.*)")

type SmtpAuthPlainMessage struct {
}

func (s SmtpAuthPlainMessage) Matches(arg string) bool {
	return authPlainRegex.MatchString(arg)
}

func (s SmtpAuthPlainMessage) Handle(transaction *SmtpTransaction, arg string) string {
	matches := authPlainRegex.FindStringSubmatch(arg)
	if len(matches) >= 2 {
		transaction.B64PlainAuth = matches[1]

		decoded, err := base64.StdEncoding.DecodeString(matches[1])
		if err == nil {
			parts := bytes.Split(decoded, []byte{0})
			if len(parts) == 3 {
				transaction.AuthorizationId = string(parts[0])
				transaction.Username = string(parts[1])
				transaction.Password = string(parts[2])
			}
		}
	} else {
		fmt.Println("warning: plain auth matched but didn't have submatch")
	}

	return "235 2.7.0 Authentication successful"
}
