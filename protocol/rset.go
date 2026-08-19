package protocol

type SmtpRsetMessage struct {
}

func (s SmtpRsetMessage) Matches(arg []byte) bool {
	return string(arg) == "RSET"
}

func (s SmtpRsetMessage) Handle(connection *SmtpConnection, arg []byte) string {
	curTrans := connection.GetCurrentTransaction(false)
	curTrans.Status = SMTP_TRANSACTION_STATUS_RESET

	return "250 2.0.0 OK"
}
