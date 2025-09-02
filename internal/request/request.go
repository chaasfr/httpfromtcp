package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"
)

const crlf = "\r\n"

const bufferSize = 8

type RequestState int

const (
	INITIALIZED RequestState = iota
	DONE
)

type Request struct {
	State       RequestState
	RequestLine RequestLine
	// Headers     map[string]string
	// Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case INITIALIZED:
		nbrBytesConsumed, requestLine, err := parseRequestLine(data)
		if err != nil {
			return -1, err
		}
		if nbrBytesConsumed == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.State = DONE
		return nbrBytesConsumed, nil
	case DONE:
		return -1, errors.New("error: trying to read data in a done state")
	default:
		return -1, errors.New("error: unknown state")
	}
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil, nil
	}

	text := string(data[:idx])
	sliced := strings.Split(text, "\r\n")
	requestLinesplit := strings.Split(sliced[0], " ")

	if len(requestLinesplit) != 3 {
		return idx, nil, errors.New("Wrong number of arguments in requestLine: " + text)
	}

	requestLine := RequestLine{
		HttpVersion: strings.Split(requestLinesplit[2],"/")[1],
		RequestTarget: requestLinesplit[1],
		Method: requestLinesplit[0],
	}

	if requestLine.HttpVersion != "1.1" {
		return idx, nil, errors.New("wrong http version "+ requestLine.HttpVersion + " : we only support 1.1")
	}

	for _, l := range requestLine.Method {
		if !unicode.IsLetter(l) || !unicode.IsUpper(l) {
			return idx, nil, errors.New("wrong method " + requestLine.Method + " : must be only capitalized letters")
		}
	}

	return idx, &requestLine, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data := make([]byte,bufferSize, bufferSize)
	readToIndex := 0
	req := Request{
		State: INITIALIZED,
	}

	for req.State != DONE {
		if readToIndex == len(data) {
			newData := make([]byte, len(data)*2, len(data)*2)
			copy(newData, data)
			data = newData
		}
		nbrBytesRead, err := reader.Read(data[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.State = DONE
				break
			}
			return nil, err
		}
		readToIndex += nbrBytesRead
		nbrBytesParsed, err := req.parse(data)
		if err != nil {
			return nil, err
		}
		newData := make([]byte, len(data)-nbrBytesParsed,len(data)-nbrBytesParsed)
		copy(newData, data[nbrBytesParsed:])
		readToIndex -= nbrBytesParsed
	}

	return &req, nil
}