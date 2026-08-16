package protocol

import (
	"fmt"
	"regexp"
)

var etrnRegex = regexp.MustCompile("ETRN (.*)")

type SmtpEtrnMessage struct {
}

func (s SmtpEtrnMessage) Matches(arg []byte) bool {
	return etrnRegex.Match(arg)
}

func (s SmtpEtrnMessage) Handle(transaction *SmtpTransaction, arg []byte) string {
	matches := etrnRegex.FindSubmatch(arg)
	if len(matches) >= 2 {
		return fmt.Sprintf("250 Queuing for node %s started", string(matches[1]))
	}

	fmt.Println("warning: etrn missing submatch")
	return "250 Queuing started"
}
