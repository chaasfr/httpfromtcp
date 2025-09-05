package main

import (
	"HTTPFROMTCP/internal/request"
	"HTTPFROMTCP/internal/response"
	"HTTPFROMTCP/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	hErr := server.HandlerError{
		StatusCode : response.OK,
		Title: "200 OK",
		SubTitle: "Success!",
		Message: "Your request was an absolute banger.",
	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		hErr.StatusCode = response.BAD_REQ
		hErr.Title = "400 Bad Request"
		hErr.SubTitle = "Bad Request"
		hErr.Message = "Your request honestly kinda sucked."
	case "/myproblem":
		hErr.StatusCode = response.INTERNAL_ERROR
		hErr.Title = "500 Internal Server Error"
		hErr.SubTitle = "Internal Server Error"
		hErr.Message = "Okay, you know what? This one is on me."
	}

	body := []byte(hErr.ToHTML())
	h := response.GetDefaultHeaders(len(body))
	h.Replace("content-type", "text/html")

	w.WriteStatusLine(hErr.StatusCode)
	w.WriteHeaders(h)
	w.WriteBody(body)
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