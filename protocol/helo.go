package protocol

import (
	"regexp"
)

var heloRegex = regexp.MustCompile("(?i)(HELO$)|(HELO (.*))")

type SmtpHeloMessage struct{}

func (s SmtpHeloMessage) Matches(arg []byte) bool {
	return heloRegex.Match(arg)
}

func (s SmtpHeloMessage) Handle(connection *SmtpConnection, arg []byte) string {
	matches := heloRegex.FindSubmatch(arg)
	if len(matches) >= 4 {
		connection.Hostname = string(matches[3])
	}

	return "250 example.com"
}
