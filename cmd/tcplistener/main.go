package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(conn net.Conn) <-chan string {
	result := make(chan string)

	go func() {
		defer conn.Close()
		defer close(result)
		newLine := ""
		for {
			slice := make([]byte, 8)
			_, err := conn.Read(slice)
			if err != nil {
				if newLine != "" {
					result <- newLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			sliced := strings.Split(string(slice), "\n")
			for i := 0 ; i < len(sliced)-1; i++ {
				result <- fmt.Sprintf("%s%s", newLine, sliced[i])
				newLine = ""
			}
			newLine += sliced[len(sliced)-1]
		}
	}()

	return result
}

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
		linesChannel := getLinesChannel(conn)
		for line := range linesChannel {
			fmt.Println(line)
		}
		fmt.Println("==== Connection ", conn.RemoteAddr()," has been closed ====")
	}
}