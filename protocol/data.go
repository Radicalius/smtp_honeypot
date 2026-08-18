package protocol

import (
	"encoding/base64"
	"fmt"
)

type SmtpDataMessage struct{}

func (s SmtpDataMessage) Matches(message []byte) bool {
	return string(message) == "DATA"
}

func (s SmtpDataMessage) Handle(connection *SmtpConnection, message []byte) string {
	if s.Matches(message) {
		connection.Deferred = s
		return "354 3.0.0 Start mail input"
	}

	lastTransaction := connection.GetCurrentTransaction()

	if string(message) != "." {
		lastTransaction.RawData = append(lastTransaction.RawData, message...)
		return ""
	}

	lastTransaction.B64Data = base64.StdEncoding.EncodeToString(lastTransaction.RawData)
	lastTransaction.Status = SMTP_TRANSACTION_STATUS_COMPLETE
	connection.Deferred = nil
	fmt.Printf("finished transation with contents: %s\n", lastTransaction.B64Data)
	return "250 2.0.0 OK"
}
