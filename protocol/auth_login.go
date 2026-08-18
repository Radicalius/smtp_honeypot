package protocol

import "encoding/base64"

type SmtpAuthLoginMessage struct {
}

func (s SmtpAuthLoginMessage) Matches(arg []byte) bool {
	return string(arg) == "AUTH LOGIN"
}

func (s SmtpAuthLoginMessage) Handle(connection *SmtpConnection, arg []byte) string {
	if string(arg) == "AUTH LOGIN" {
		connection.Authentication = append(connection.Authentication, SmtpAuthentication{
			Type: "LOGIN",
		})
		connection.Deferred = s
		return "334 VXNlcm5hbWU6"
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
