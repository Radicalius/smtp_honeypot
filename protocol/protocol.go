package protocol

type SmtpTransaction struct {
	Hostname         string
	From             []string
	To               []string
	Data             string
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
	Matches(string) bool
	Handle(*SmtpTransaction, string) string
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
		return transaction.Deferred.Handle(transaction, string(body))
	}

	for _, message := range smtpMessages {
		if message.Matches(string(body)) {
			return message.Handle(transaction, string(body))
		}
	}

	return "503 5.5.1 Bad sequence"
}
