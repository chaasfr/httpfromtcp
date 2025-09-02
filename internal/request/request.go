package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	// Headers     map[string]string
	// Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(data string) (*Request, error) {
	sliced := strings.Split(data, "\r\n")
	requestLinesplit := strings.Split(sliced[0], " ")

	if len(requestLinesplit) != 3 {
		return nil, errors.New("Wrong number of arguments in requestLine: " + data)
	}

	requestLine := RequestLine{
		HttpVersion: strings.Split(requestLinesplit[2],"/")[1],
		RequestTarget: requestLinesplit[1],
		Method: requestLinesplit[0],
	}

	if requestLine.HttpVersion != "1.1" {
		return nil, errors.New("wrong http version "+ requestLine.HttpVersion + " : we only support 1.1")
	}

	for _, l := range requestLine.Method {
		if !unicode.IsLetter(l) || !unicode.IsUpper(l) {
			return nil, errors.New("wrong method " + requestLine.Method + " : must be only capitalized letters")
		}
	}

	request := Request{
		RequestLine: requestLine,
	}

	return &request, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	req, err := parseRequestLine(string(data[:]))
	if err != nil {
		return nil, err
	}

	return req, nil
}