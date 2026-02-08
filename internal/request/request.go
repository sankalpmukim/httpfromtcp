package request

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sankalpmukim/httpfromtcp/internal/headers"
)

type RequestState int

const (
	requestStateInitialized RequestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
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
	n := 0

	req := &Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	for req.state != requestStateDone {

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

		for consumed > 0 {
			consumed, err = req.parse(buf[:readToIndex])
			if err != nil {
				return nil, err
			}

			if consumed > 0 {
				// slide remaining bytes to front
				copy(buf, buf[consumed:readToIndex])
				readToIndex -= consumed
			}
		}

		// Grow if full
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		if req.state == requestStateDone {
			break
		}

		n, err = reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {

				// case where it EOF'd but *if* unparsed buffer exists
				// so we still need to parse it.
				if n != readToIndex || readToIndex != 0 {
					_, err := req.parse(buf[:readToIndex])
					if err != nil {
						return nil, err
					}
				}

				if req.state == requestStateDone {
					break
				} else if req.state == requestStateParsingBody {
					// case where req.state == body parsing, but headers suggest body not expected.
					// i.e., expected the tcp connection to get closed.
					contentLengthStr := req.Headers.Get("Content-Length")
					if contentLengthStr == "" {
						req.state = requestStateDone
						continue
					} else {
						contentLength, err := strconv.Atoi(contentLengthStr)
						if err != nil {
							return nil, err
						}
						if contentLength == 0 {
							req.state = requestStateDone
							continue
						}
						return req, errors.New("Connection ended abruptly, before headers ended")
					}
				} else {
					return req, errors.New("Connection ended abruptly, before headers ended")
				}
			}
			return nil, err
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
		bytesRead, headersDone, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if headersDone {
			r.state = requestStateParsingBody
		}
		return bytesRead, nil

	case requestStateParsingBody:
		contentLengthStr := r.Headers.Get("Content-Length")
		if contentLengthStr == "" {
			r.state = requestStateDone
			return 0, nil
		} else {
			contentLength, err := strconv.Atoi(contentLengthStr)
			if err != nil {
				return 0, err
			}
			r.Body = append(r.Body, data...)
			if contentLength < len(r.Body) {
				return 0, fmt.Errorf("Received more body than Content-Length provided. Content-Length=%v, consumed body=%v", contentLength, len(r.Body))
			}
			// if contentLength > len(r.Body) {
			// 	return 0, fmt.Errorf("Content-Length greater than consumed body length")
			// }
			if contentLength == len(r.Body) {
				r.state = requestStateDone
			}
			return len(data), nil
		}

	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")

	default:
		return 0, errors.New("error: unknown state")
	}
}
