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

func (s SmtpAuthPlainMessage) Matches(arg []byte) bool {
	return authPlainRegex.Match(arg)
}

func (s SmtpAuthPlainMessage) Handle(connection *SmtpConnection, arg []byte) string {
	matches := authPlainRegex.FindSubmatch(arg)

	auth := SmtpAuthentication{
		Type: "PLAIN",
	}
	if len(matches) >= 2 {
		auth.B64PlainAuth = string(matches[1])

		decoded, err := base64.StdEncoding.DecodeString(string(matches[1]))
		if err == nil {
			parts := bytes.Split(decoded, []byte{0})
			if len(parts) == 3 {
				auth.AuthorizationId = string(parts[0])
				auth.Username = string(parts[1])
				auth.Password = string(parts[2])
			}
		}
	} else {
		fmt.Println("warning: plain auth matched but didn't have submatch")
	}

	connection.Authentication = append(connection.Authentication, auth)

	return "235 2.7.0 Authentication successful"
}
