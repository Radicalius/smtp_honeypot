package protocol

import "encoding/base64"

type SmtpAuthLoginMessage struct {
}

func (s SmtpAuthLoginMessage) Matches(arg []byte) bool {
	return string(arg) == "AUTH LOGIN"
}

func (s SmtpAuthLoginMessage) Handle(transaction *SmtpTransaction, arg []byte) string {
	if string(arg) == "AUTH LOGIN" {
		transaction.Username = ""
		transaction.Password = ""
		transaction.Deferred = s
		return "334 VXNlcm5hbWU6"
	}

	if transaction.Username == "" {
		transaction.B64Username = string(arg)
		decoded, err := base64.StdEncoding.DecodeString(string(arg))
		if err == nil {
			transaction.Username = string(decoded)
		}

		return "334 UGFzc3dvcmQ6"
	}

	transaction.B64Password = string(arg)
	decoded, err := base64.StdEncoding.DecodeString(string(arg))
	if err == nil {
		transaction.Password = string(decoded)
	}

	transaction.Deferred = nil
	return "235 2.7.0 Authentication successful"
}
