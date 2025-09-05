package server

import (
	"HTTPFROMTCP/internal/request"
	"HTTPFROMTCP/internal/response"
	"fmt"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Title      string
	SubTitle   string
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)


func (he HandlerError) Write(w *response.Writer) {
	w.WriteStatusLine(he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	w.WriteHeaders(headers)
	w.WriteBody(messageBytes)
}

func (he HandlerError) ToHTML() string {
	return fmt.Sprintf(
		`<html>
  <head>
    <title>%s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>%s</p>
  </body>
</html>`, he.Title, he.SubTitle, he.Message)
}