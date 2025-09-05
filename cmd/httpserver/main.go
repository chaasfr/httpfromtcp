package main

import (
	"HTTPFROMTCP/internal/request"
	"HTTPFROMTCP/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	msg := ""
	statusCode := 0
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		statusCode = 400
		msg = "Your problem is not my problem\n"
	case "/myproblem":
		statusCode = 500
		msg = "Woopsie, my bad\n"
	default:
		statusCode = 200
		msg = "All good, frfr\n"
	}
	_, err := w.Write([]byte(msg))
	if err != nil {
		return &server.HandlerError{
			StatusCode: 500,
			Message: err.Error(),
		}
	}
	return &server.HandlerError{
		StatusCode: statusCode,
		Message: msg,
	}
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}