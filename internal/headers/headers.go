package headers

import (
	"bytes"
	"errors"
	"strings"
)

const crlf = "\r\n"

const allowedCharsForKey="abcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"

type Headers map[string]string

func NewHeaders()(Headers) {
	return Headers{}
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
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
		return 0, false, errors.New("wrong header format: missing a ':' in "+ text)
	}

	if strings.HasSuffix(sliced[0], " ") {
		return 0, false, errors.New("wrong header format: unexpected whitespace before ':'")
	}

	key := strings.ToLower(strings.TrimLeft(sliced[0], " "))

	if len(key) < 1 {
		return 0, false, errors.New("header key must have a length of 1 or more")
	}

	for _, l := range key {
		if !strings.Contains(allowedCharsForKey,string(l)) {
			return 0, false, errors.New("wrong character " + string(l) + "in header key " + key)
		}
	}

	value := strings.Trim(sliced[1], " ")

	h.Set(key, string(value))
	return idx + 2, false, nil
}
