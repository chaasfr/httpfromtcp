package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func getLinesChannel(conn net.Conn) <-chan string {
	result := make(chan string)

	go func() {
		defer conn.Close()
		defer close(result)
		var newLine string
		for {
			slice := make([]byte, 8)
			_, err := conn.Read(slice)
			if err == io.EOF{
				break
			}
			sliced := strings.Split(string(slice), "\n")
			newLine += sliced[0]
			for i := 1 ; i < len(sliced)-1; i++ {
				result <- newLine
				newLine = sliced[i]
			}
		}

		if newLine != "" {
			result <- newLine
		}
	}()

	return result
}

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("could not listen %s: %s\n", ln, err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
		log.Fatalf("could not open conn %s: %s\n", conn, err)
		os.Exit(1)
		}
		fmt.Println("==== Connection has been accepted ====")
		linesChannel := getLinesChannel(conn)
		for line := range linesChannel {
			fmt.Printf("read: %s\n", line)
		}
		fmt.Println("==== Connection has been closed ====")
	}
}