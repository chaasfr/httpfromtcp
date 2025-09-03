package headers

import (
	"bytes"
	"errors"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders()(Headers) {
	return Headers{}
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	
	if idx == 0 {
		return 2, true, nil
	}

	text := string(data[:idx])
	sliced := strings.SplitN(text, ":",2)

	if len(sliced) < 2 {
		return 0, false, errors.New("wrong header format: missing a ':'")
	}

	if strings.HasSuffix(sliced[0], " ") {
		return 0, false, errors.New("wrong header format: unexpected whitespace before ':'")
	}

	key := strings.TrimLeft(sliced[0], " ")
	value := strings.Trim(sliced[1], " ")

	h.Set(key, string(value))
	return idx + 2, false, nil
}
