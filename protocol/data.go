package protocol

import "encoding/base64"

type SmtpDataMessage struct{}

func (s SmtpDataMessage) Matches(message string) bool {
	return message == "DATA"
}

func (s SmtpDataMessage) Handle(transaction *SmtpTransaction, message string) string {
	if s.Matches(message) {
		transaction.Deferred = s
		return "354 3.0.0 Start mail input"
	}

	transaction.Data = base64.StdEncoding.EncodeToString([]byte(message))
	transaction.Deferred = nil
	return "250 2.0.0 OK"
}
