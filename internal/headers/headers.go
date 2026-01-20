package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

// consumes all headers at once, and stores them in the Headers object.
// if data does not contain CRLF, it returns early as
// it does not have enough data yet
func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	allContent := string(data)
	idx := strings.Index(allContent, "\r\n")
	if idx == -1 {
		return 0, false, nil
	}

	lineContent := strings.TrimSpace(allContent[:idx])
	bytesConsumed := idx + 2

	if len(lineContent) == 0 {
		return 0, true, nil
	}

	idxColon := strings.Index(lineContent, ":")

	if idxColon <= 0 || lineContent[idxColon-1] == ' ' {
		return 0, false, errors.ErrUnsupported
	}

	h[strings.TrimSpace(lineContent[:idxColon])] = strings.TrimSpace(lineContent[idxColon+2:])

	return bytesConsumed, false, nil
}

func NewHeaders() Headers {
	headers := make(map[string]string)
	return headers
}
