package protocol

import (
	"bufio"
	"fmt"
	"regexp"
)

var etrnRegex = regexp.MustCompile("(?i)ETRN (.*)")

type SmtpEtrnMessage struct {
}

func (s SmtpEtrnMessage) Matches(arg []byte) bool {
	return etrnRegex.Match(arg)
}

func (s SmtpEtrnMessage) Handle(connection *SmtpConnection, arg []byte, _ *bufio.Reader) string {
	connection.EtrnEnabled = true

	matches := etrnRegex.FindSubmatch(arg)
	if len(matches) >= 2 {
		connection.EtrnNode = string(matches[1])
		return fmt.Sprintf("250 Queuing for node %s started", connection.EtrnNode)
	}

	fmt.Println("warning: etrn missing submatch")
	return "250 Queuing started"
}
