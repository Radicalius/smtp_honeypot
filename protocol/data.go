package protocol

import "encoding/base64"

type SmtpDataMessage struct{}

func (s SmtpDataMessage) Matches(message []byte) bool {
	return string(message) == "DATA"
}

func (s SmtpDataMessage) Handle(transaction *SmtpTransaction, message []byte) string {
	if s.Matches(message) {
		transaction.Deferred = s
		return "354 3.0.0 Start mail input"
	}

	if string(message) != "." {
		transaction.RawData = append(transaction.RawData, message...)
		return ""
	}

	transaction.B64Data = base64.StdEncoding.EncodeToString(transaction.RawData)
	transaction.Deferred = nil
	return "250 2.0.0 OK"
}
