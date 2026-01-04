package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(reqContent string) (RequestLine, error) {
	requestLineStr := strings.Split(reqContent, "\r\n")[0]
	reqLineParts := strings.Split(requestLineStr, " ")
	if len(reqLineParts) != 3 {
		return RequestLine{}, errors.New("invalid request line format")
	}
	method := reqLineParts[0]
	requestTarget := reqLineParts[1]
	httpVersion := strings.Split(reqLineParts[2], "/")[1]
	return RequestLine{httpVersion, requestTarget, method}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqStringBytes, err := io.ReadAll(reader)
	reqString := string(reqStringBytes)
	if err != nil {
		return nil, fmt.Errorf("Error in io.ReadAll %v", err)
	}
	parsedReqLine, err := parseRequestLine(reqString)
	if err != nil {
		return nil, err
	}
	return &Request{parsedReqLine}, nil
}
