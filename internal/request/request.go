package request

import (
	"HTTPFROMTCP/internal/headers"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const crlf = "\r\n"

const bufferSize = 8

type RequestState int

const (
	INITIALIZED RequestState = iota
	PARSING_HEADERS
	PARSING_BODY
	DONE
)

type Request struct {
	state       RequestState
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case INITIALIZED:
		nbrBytesConsumed, requestLine, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if nbrBytesConsumed == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = PARSING_HEADERS
		return nbrBytesConsumed, nil
	case PARSING_HEADERS:
		nbrBytesConsumed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = PARSING_BODY
		}
		return nbrBytesConsumed, nil
	case PARSING_BODY:
		contentLength := r.Headers.Get("content-length")
		if contentLength == "" {
			r.state = DONE
			return 0, nil
		}
		contentLengthInt, err := strconv.Atoi(contentLength)
		if err != nil {
			return 0, err
		}
		r.Body = append(r.Body,data...)
		if int(contentLengthInt) < len(r.Body) {
			return 0, errors.New("body is longer than contentLength of " +
			 contentLength +
			" :" + string(r.Body))
		}
		if int(contentLengthInt) == len(r.Body) {
			r.state = DONE
		}
		return len(data), nil
	case DONE:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != DONE {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil, nil
	}

	text := string(data[:idx])
	requestLinesplit := strings.Split(text, " ")

	if len(requestLinesplit) != 3 {
		return idx, nil, errors.New("Wrong number of arguments in requestLine: " + text)
	}

	requestLine := RequestLine{
		HttpVersion:   strings.Split(requestLinesplit[2], "/")[1],
		RequestTarget: requestLinesplit[1],
		Method:        requestLinesplit[0],
	}

	if requestLine.HttpVersion != "1.1" {
		return idx, nil, errors.New("wrong http version " + requestLine.HttpVersion + " : we only support 1.1")
	}

	for _, l := range requestLine.Method {
		if !unicode.IsLetter(l) || !unicode.IsUpper(l) {
			return idx, nil, errors.New("wrong method " + requestLine.Method + " : must be only capitalized letters")
		}
	}

	return idx + 2, &requestLine, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		state:   INITIALIZED,
		Headers: headers.NewHeaders(),
	}

	for req.state != DONE {
		if readToIndex >= len(data) {
			newData := make([]byte, len(data)*2)
			copy(newData, data)
			data = newData
		}
		nbrBytesRead, err := reader.Read(data[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != DONE {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.state, nbrBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += nbrBytesRead		
		nbrBytesParsed, err := req.parse(data[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(data, data[nbrBytesParsed:])
		readToIndex -= nbrBytesParsed
	}

	return req, nil
}
