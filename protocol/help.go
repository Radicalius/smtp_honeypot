package protocol

import (
	"bytes"
	"regexp"
	"strings"
)

type SmtpHelpMessage struct{}

var helpRegex = regexp.MustCompile("(HELP$)|(HELP (.*))")

func (s SmtpHelpMessage) Matches(arg []byte) bool {
	return bytes.Contains(arg, []byte("HELP"))
}

func (s SmtpHelpMessage) Handle(connection *SmtpConnection, arg []byte) string {
	matches := helpRegex.FindSubmatch(arg)
	if len(matches) >= 4 && len(matches[3]) > 0 {
		return "214 Help entry for " + string(matches[3])
	}

	return strings.ReplaceAll(`214-This server supports the following commands:
214-HELO EHLO MAIL RCPT DATA RSET VRFY EXPN HELP NOOP QUIT STARTTLS AUTH
214 End of HELP info`, "\n", "\r\n")
}
