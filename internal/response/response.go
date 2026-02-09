package response

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sankalpmukim/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := fmt.Sprintf("HTTP/1.1 %v ", statusCode)
	switch statusCode {
	case OK:
		statusLine += "OK"
	case BadRequest:
		statusLine += "Bad Request"
	case InternalServerError:
		statusLine += "Internal Server Error"
	default:
		statusLine += ""
	}
	_, err := w.Write([]byte(statusLine + "\r\n"))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = strconv.Itoa(contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "plain"
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var headersString strings.Builder
	for k, v := range headers {
		fmt.Fprintf(&headersString, "%v: %v\r\n", k, v)
	}
	fmt.Fprint(&headersString, "\r\n")
	_, err := w.Write([]byte(headersString.String()))
	return err
}

type writerState int

const (
	statusLineNext writerState = iota
	headersNext
	bodyNext
)

type Writer struct {
	io.Writer
	toWriteNext writerState
}

func NewResponseWriter(w io.Writer) Writer {
	return Writer{Writer: w, toWriteNext: statusLineNext}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.toWriteNext == statusLineNext {
		w.toWriteNext = headersNext
		return WriteStatusLine(w, statusCode)
	} else {
		return errors.ErrUnsupported
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.toWriteNext == headersNext {
		w.toWriteNext = bodyNext
		return WriteHeaders(w, headers)
	} else {
		return errors.ErrUnsupported
	}
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.toWriteNext == bodyNext {
		return w.Write(p)
	} else {
		return 0, errors.ErrUnsupported
	}
}
