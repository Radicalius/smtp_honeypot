package protocol

import (
	"fmt"
	"regexp"
	"strings"
)

var vrfyRegex = regexp.MustCompile("VRFY (.*)")

type SmtpVrfyMessage struct {
}

func (s SmtpVrfyMessage) Matches(arg []byte) bool {
	return vrfyRegex.Match(arg)
}

func (s SmtpVrfyMessage) Handle(connection *SmtpConnection, arg []byte) string {
	matches := vrfyRegex.FindSubmatch(arg)
	if len(matches) >= 2 {
		addr := strings.Replace(string(matches[1]), "\r\n", "", 1)
		return fmt.Sprintf("250 <%s>", addr)
	}

	return "501 5.1.3 Bad recipient address syntax"
}
