package response

import (
	"HTTPFROMTCP/internal/headers"
	"errors"
	"fmt"
	"io"
)

type StatusCode int

const (
	OK             StatusCode = 200
	BAD_REQ        StatusCode = 400
	INTERNAL_ERROR StatusCode = 500
)

type WriterState int

const (
	AWAITING_STATUS_LINE WriterState = iota
	AWAITING_HEADERS
	AWAITING_BODY
	DONE
)

type Writer struct {
	StatusLine []byte
	Headers    []byte
	Body       []byte
	Writer     io.Writer
	State      WriterState
}

func NewWriter(w io.Writer) *Writer{
	return &Writer{
		State: AWAITING_STATUS_LINE,
		Writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != AWAITING_STATUS_LINE {
		return errors.New("error: trying to write status line during invalid state")
	}
	var msg string
	switch statusCode {
	case OK:
		msg = "HTTP/1.1 200 OK"
	case BAD_REQ:
		msg = "HTTP/1.1 400 Bad Request"
	case INTERNAL_ERROR:
		msg = "HTTP/1.1 500 Server Error"
	default:
		msg = ""
	}
	w.StatusLine = []byte(msg + "\r\n")
	_, err := w.Writer.Write(w.StatusLine)
	w.State = AWAITING_HEADERS
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("content-length", fmt.Sprint(contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")

	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.State != AWAITING_HEADERS {
		return errors.New("error: trying to write headers during invalid state")
	}
	for k, v := range headers {
		txt := k + ": " + v + "\r\n"
		w.Headers = append(w.Headers,[]byte(txt)...)
	}
	w.Headers = append(w.Headers, []byte("\r\n")...)

	_, err := w.Writer.Write(w.Headers)
	w.State = AWAITING_BODY
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.State != AWAITING_BODY {
		return 0, errors.New("error: trying to write body during invalid state")
	}
	w.Body = p
	return w.Writer.Write(w.Body)
}