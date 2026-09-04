package protocol

import (
	"fmt"
	"regexp"
)

var fromRegex = regexp.MustCompile("(?i)MAIL FROM:<([^>]*)>")

type SmtpFromMessage struct {
}

func (s SmtpFromMessage) Matches(arg []byte) bool {
	return fromRegex.Match(arg)
}

func (s SmtpFromMessage) Handle(connection *SmtpConnection, arg []byte) string {
	matches := fromRegex.FindSubmatch(arg)
	lastTrans := connection.GetCurrentTransaction(true)
	if len(matches) < 2 {
		fmt.Println("warning: no from match")
	} else if len(lastTrans.From) >= 100 {
		fmt.Println("warning: more than 100 from fields supplied")
	} else {
		lastTrans.From = append(lastTrans.From, string(matches[1]))
	}

	return "250 2.1.0 OK"
}
