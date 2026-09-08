package protocol

import (
	"encoding/base64"
	"regexp"
)

var authLoginRegex = regexp.MustCompile(`(?i)^AUTH LOGIN(?:\s+(\S+))?$`)

type SmtpAuthLoginMessage struct {
}

func (s SmtpAuthLoginMessage) Matches(arg []byte) bool {
	return authLoginRegex.Match(arg)
}

func (s SmtpAuthLoginMessage) Handle(connection *SmtpConnection, arg []byte) string {
	matches := authLoginRegex.FindSubmatch(arg)
	if len(matches) > 0 {
		auth := SmtpAuthentication{Type: "LOGIN"}
		if len(matches) > 1 && len(matches[1]) > 0 {
			auth.B64Username = string(matches[1])
			decoded, err := base64.StdEncoding.DecodeString(string(matches[1]))
			if err == nil {
				auth.Username = string(decoded)
			}
		}

		connection.Authentication = append(connection.Authentication, auth)
		connection.Deferred = s
		if auth.Username == "" {
			return "334 VXNlcm5hbWU6"
		}

		return "334 UGFzc3dvcmQ6"
	}

	lastAuth := &connection.Authentication[len(connection.Authentication)-1]
	if lastAuth.Username == "" {
		lastAuth.B64Username = string(arg)
		decoded, err := base64.StdEncoding.DecodeString(string(arg))
		if err == nil {
			lastAuth.Username = string(decoded)
		}

		return "334 UGFzc3dvcmQ6"
	}

	lastAuth.B64Password = string(arg)
	decoded, err := base64.StdEncoding.DecodeString(string(arg))
	if err == nil {
		lastAuth.Password = string(decoded)
	}

	connection.Deferred = nil
	return "235 2.7.0 Authentication successful"
}
