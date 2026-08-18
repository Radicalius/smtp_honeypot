package protocol

import (
	"fmt"
	"regexp"
)

var toRegex = regexp.MustCompile("RCPT TO:<([^>]*)>")

type SmtpToMessage struct {
}

func (s SmtpToMessage) Matches(arg []byte) bool {
	return toRegex.Match(arg)
}

func (s SmtpToMessage) Handle(connection *SmtpConnection, arg []byte) string {
	matches := toRegex.FindSubmatch(arg)
	lastTrans := connection.GetCurrentTransaction()
	if len(matches) < 2 {
		fmt.Println("warning: no to match")
	} else if len(lastTrans.To) >= 100 {
		fmt.Println("warning: more than 100 to fields supplied")
	} else {
		lastTrans.To = append(lastTrans.To, string(matches[1]))
	}

	return "250 2.1.5 OK"
}
