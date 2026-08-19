package protocol

type SmtpAuthentication struct {
	Type            string `json:"type"`
	B64PlainAuth    string `json:"-"`
	AuthorizationId string `json:"authorizationId"`
	Username        string `json:"username"`
	B64Username     string `json:"-"`
	Password        string `json:"password"`
	B64Password     string `json:"-"`
}

const (
	SMTP_TRANSACTION_STATUS_IN_PROGRESS = 0
	SMTP_TRANSACTION_STATUS_RESET       = 1
	SMTP_TRANSACTION_STATUS_COMPLETE    = 2
)

type SmtpTransactionStatus uint8

type SmtpTransaction struct {
	Status  SmtpTransactionStatus `json:"status"`
	From    []string              `json:"from"`
	To      []string              `json:"to"`
	RawData []byte                `json:"-"`
	B64Data string                `json:"data"`
}

type SmtpConnection struct {
	Guid             string               `json:"guid"`
	SrcAddr          string               `json:"srcAddr"`
	Hostname         string               `json:"hostname"`
	Transactions     []SmtpTransaction    `json:"transactions"`
	Authentication   []SmtpAuthentication `json:"authentication"`
	TLS              bool                 `json:"tls"`
	ExtendedProtocol bool                 `json:"extended"`
	StartEpochMs     uint64               `json:"startEpochMs"`
	DurationMs       uint64               `json:"durationMs"`

	Deferred SmtpMessage `json:"-"`
}

func (s *SmtpConnection) GetCurrentTransaction(createNewIfComplete bool) *SmtpTransaction {
	if len(s.Transactions) == 0 {
		s.Transactions = []SmtpTransaction{SmtpTransaction{}}
	}

	lastTrans := &s.Transactions[len(s.Transactions)-1]
	if lastTrans.Status != SMTP_TRANSACTION_STATUS_IN_PROGRESS && createNewIfComplete {
		s.Transactions = append(s.Transactions, SmtpTransaction{})
		lastTrans = &s.Transactions[len(s.Transactions)-1]
	}

	return lastTrans
}

type SmtpMessage interface {
	Matches([]byte) bool
	Handle(*SmtpConnection, []byte) string
}

var smtpMessages []SmtpMessage = []SmtpMessage{
	SmtpHeloMessage{},
	SmtpEhloMessage{},
	SmtpFromMessage{},
	SmtpToMessage{},
	SmtpDataMessage{},
	SmtpAuthLoginMessage{},
	SmtpAuthPlainMessage{},
	SmtpVrfyMessage{},
	SmtpEtrnMessage{},
	SmtpRsetMessage{},
}

func Handle(connection *SmtpConnection, body []byte) string {
	if connection.Deferred != nil {
		return connection.Deferred.Handle(connection, body)
	}

	for _, message := range smtpMessages {
		if message.Matches(body) {
			return message.Handle(connection, body)
		}
	}

	return "503 5.5.1 Bad sequence"
}
