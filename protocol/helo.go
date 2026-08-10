package protocol

import "regexp"

var heloRegex = regexp.MustCompile("HELO (.*)")

type SmtpHeloMessage struct{}

func (s SmtpHeloMessage) Matches(arg string) bool {
	return heloRegex.MatchString(arg)
}

func (s SmtpHeloMessage) Handle(transaction *SmtpTransaction, arg string) string {
	matches := heloRegex.FindStringSubmatch(arg)
	transaction.Hostname = matches[1]
	return "250 example.com"
}
