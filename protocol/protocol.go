package protocol

type SmtpTransaction struct {
	Hostname         string
	From             []string
	To               []string
	RawData          []byte
	B64Data          string
	B64PlainAuth     string
	AuthorizationId  string
	Username         string
	B64Username      string
	Password         string
	B64Password      string
	ExtendedProtocol bool
	Deferred         SmtpMessage
}

type SmtpMessage interface {
	Matches([]byte) bool
	Handle(*SmtpTransaction, []byte) string
}

var smtpMessages []SmtpMessage = []SmtpMessage{
	SmtpHeloMessage{},
	SmtpEhloMessage{},
	SmtpFromMessage{},
	SmtpToMessage{},
	SmtpDataMessage{},
	SmtpAuthLoginMessage{},
	SmtpAuthPlainMessage{},
}

func Handle(transaction *SmtpTransaction, body []byte) string {
	if transaction.Deferred != nil {
		return transaction.Deferred.Handle(transaction, body)
	}

	for _, message := range smtpMessages {
		if message.Matches(body) {
			return message.Handle(transaction, body)
		}
	}

	return "503 5.5.1 Bad sequence"
}
