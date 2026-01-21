package request

import (
	"errors"
	"io"
	"strings"

	"github.com/sankalpmukim/httpfromtcp/internal/headers"
)

type RequestState int

const (
	requestStateInitialized RequestState = iota
	requestStateParsingHeaders
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       RequestState // 0 -> initialized, 1 -> done
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const bufferSize = 8

func parseRequestLine(data string) (int, RequestLine, error) {
	// Find end of request line
	idx := strings.Index(data, "\r\n")
	if idx == -1 {
		// Need more data
		return 0, RequestLine{}, nil
	}

	line := data[:idx] // without \r\n
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return 0, RequestLine{}, errors.New("invalid request line format")
	}

	method := parts[0]
	target := parts[1]

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 {
		return 0, RequestLine{}, errors.New("invalid HTTP version")
	}
	version := httpParts[1]

	// +2 to consume the "\r\n"
	consumed := idx + 2

	return consumed, RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}, nil
}

// parseRequestLine re-implement

// func RequestFromReader(reader io.Reader) (*Request, error) {
// 	reqStringBytes, err := io.ReadAll(reader)
// 	reqString := string(reqStringBytes)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error in io.ReadAll %v", err)
// 	}
// 	parsedReqLine, err := parseRequestLine(reqString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Request{parsedReqLine}, nil
// }

// this function is called once per request, with a reader that
// can send information in chunks
func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	req := &Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	for req.state != requestStateDone {

		// Grow if full
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state == requestStateDone {
					break
				} else {
					return req, errors.New("Connection ended abruptly, before headers ended")
				}
			}
			return nil, err
		}

		readToIndex += n

		consumed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		if consumed > 0 {
			// slide remaining bytes to front
			copy(buf, buf[consumed:readToIndex])
			readToIndex -= consumed
		}
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		bytesRead, requestLine, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if bytesRead == 0 {
			return 0, nil // need more data
		}
		r.RequestLine = requestLine
		r.state = requestStateParsingHeaders
		return bytesRead, nil
	case requestStateParsingHeaders:
		bytesRead, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateDone
		}
		return bytesRead, nil

	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")

	default:
		return 0, errors.New("error: unknown state")
	}
}
