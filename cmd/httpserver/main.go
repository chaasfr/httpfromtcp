package main

import (
	"HTTPFROMTCP/internal/request"
	"HTTPFROMTCP/internal/response"
	"HTTPFROMTCP/internal/server"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handleYourProblem(w)
	case "/myproblem":
		handleMyProblem(w)
	default:
		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			handleHttpbinProx(w, req.RequestLine.RequestTarget)
		} else {
			handleSucess(w)
		}
	}
}

func handleHttpbinProx(w *response.Writer, target string) {
	endpoint := "https://httpbin.org/" + strings.TrimPrefix(target, "/httpbin/")
	fmt.Println("rerouting to "+ endpoint)
	
	resp, err := http.Get(endpoint)
	if err != nil {
		hErr := server.HandlerError{
			StatusCode : response.BAD_REQ,
			Title: "400 Bad Request",
			SubTitle: "Could not reach " + target,
			Message: err.Error(),
		}
		hErr.Write(w)
		return
	}
	defer resp.Body.Close()
	
	h := response.GetDefaultHeaders(0)
	h.Remove("content-length")
	h.Set("Transfer-Encoding", "chunked")
	w.WriteStatusLine(response.OK)
	w.WriteHeaders(h)

	n := -1
	bufferSize := 1024
	buffer := make([]byte,bufferSize)
	for err != io.EOF {
		n, err = resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			hErr := server.HandlerError{
				StatusCode : response.BAD_REQ,
				Title: "500 Intenral Error",
				SubTitle: "Could not read reply from hhtpbin",
				Message: err.Error(),
			}
			w.WriteChunkedBodyDone()
			hErr.Write(w)
			resp.Body.Close()
			return
		}
		fmt.Printf("read %d bytes from httpbin\n", n)
		_, err := w.WriteChunkedBody(buffer[:n])
		if err != nil {
			fmt.Println(err)
			hErr := server.HandlerError{
				StatusCode : response.BAD_REQ,
				Title: "500 Intenral Error",
				SubTitle: "Could not write chunked body",
				Message: err.Error(),
			}
			w.WriteChunkedBodyDone()
			hErr.Write(w)
			resp.Body.Close()
			return
		}
	}
}

func handleSucess(w *response.Writer) {
	hErr := server.HandlerError{
		StatusCode : response.OK,
		Title: "200 OK",
		SubTitle: "Success!",
		Message: "Your request was an absolute banger.",
	}
	hErr.WriteHTML(w)
}

func handleYourProblem(w *response.Writer) {
	hErr := server.HandlerError{
		StatusCode : response.BAD_REQ,
		Title: "400 Bad Request",
		SubTitle: "Bad Request",
		Message: "Your request honestly kinda sucked.",
	}
	hErr.WriteHTML(w)
}

func handleMyProblem(w *response.Writer) {
	hErr := server.HandlerError{
		StatusCode : response.INTERNAL_ERROR,
		Title: "500 Internal Server Error",
		SubTitle: "Internal Server Error",
		Message: "Okay, you know what? This one is on me.",
	}
	hErr.WriteHTML(w)
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