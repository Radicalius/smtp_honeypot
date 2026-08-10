package protocol

import "regexp"

var fromRegex = regexp.MustCompile("MAIL FROM:<([^>]*)>")

type SmtpFromMessage struct {
}

func (s SmtpFromMessage) Matches(arg string) bool {
	return fromRegex.MatchString(arg)
}

func (s SmtpFromMessage) Handle(transaction *SmtpTransaction, arg string) string {
	matches := fromRegex.FindStringSubmatch(arg)
	transaction.From = append(transaction.From, matches[1])
	return "250 2.1.0 OK"
}
