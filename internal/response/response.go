package response

import (
	"HTTPFROMTCP/internal/headers"
	"io"
	"fmt"
)

type StatusCode int

const (
	OK = 200
	BAD_REQ = 400
	INTERNAL_ERROR = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var msg string
	switch statusCode {
	case OK:
		msg = "HTTP/1.1 200 OK"
	case BAD_REQ:
		msg = "HTTP/1.1 400 Bad Request"
	case INTERNAL_ERROR:
		msg = "HTTP/1.1  500 Server Error"
	default:
		msg = ""
	}
	_, err := w.Write([]byte(msg + "\r\n"))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("content-length", fmt.Sprint(contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		txt := k + ": " + v + "\r\n"
		_, err := w.Write([]byte(txt))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}