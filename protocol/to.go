package protocol

import "regexp"

var toRegex = regexp.MustCompile("RCPT TO:<([^>]*)>")

type SmtpToMessage struct {
}

func (s SmtpToMessage) Matches(arg string) bool {
	return toRegex.MatchString(arg)
}

func (s SmtpToMessage) Handle(transaction *SmtpTransaction, arg string) string {
	matches := toRegex.FindStringSubmatch(arg)
	transaction.To = append(transaction.To, matches[1])
	return "250 2.1.5 OK"
}
