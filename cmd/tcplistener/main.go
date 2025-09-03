package main

import (
	"fmt"
	"log"
	"net"
	"HTTPFROMTCP/internal/request"
)

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("could not listen %s: %s\n", ln, err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("could not open conn %s: %s\n", conn, err)
		}
		fmt.Println("==== Connection has been accepted ====")
		req,err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("could read request: %s\n", err)
		}
		fmt.Println("Request line:")
		fmt.Println(" - Method: " + req.RequestLine.Method)
		fmt.Println(" - Target: " + req.RequestLine.RequestTarget)
		fmt.Println(" - Version: " + req.RequestLine.HttpVersion)

		fmt.Println("==== Connection ", conn.RemoteAddr()," has been closed ====")
	}
}