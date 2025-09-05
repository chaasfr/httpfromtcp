package server

import (
	"HTTPFROMTCP/internal/request"
	"fmt"
	"io"
	"log"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError


func (he HandlerError) write(w io.Writer) {
	txt := "error code: " + fmt.Sprint(he.StatusCode) + "\n messsage: " + he.Message
	_, err := w.Write([]byte(txt))
	if err != nil {
		log.Fatalf("error writing handlerError: %s\n", err)
	}
}