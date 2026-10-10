package protocol

import (
	"bufio"
	"strings"
)

type SmtpRsetMessage struct {
}

func (s SmtpRsetMessage) Matches(arg []byte) bool {
	return strings.ToUpper(string(arg)) == "RSET"
}

func (s SmtpRsetMessage) Handle(connection *SmtpConnection, arg []byte, _ *bufio.Reader) string {
	curTrans := connection.GetCurrentTransaction(false)
	curTrans.Status = SMTP_TRANSACTION_STATUS_RESET

	return "250 2.0.0 OK"
}
