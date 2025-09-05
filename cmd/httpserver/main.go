package main

import (
	"HTTPFROMTCP/internal/request"
	"HTTPFROMTCP/internal/response"
	"HTTPFROMTCP/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode : response.BAD_REQ,
			Message: "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode : response.INTERNAL_ERROR,
			Message: "Woopsie, my bad\n",
		}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
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