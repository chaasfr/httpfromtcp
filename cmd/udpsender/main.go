package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	addr, err := net.ResolveUDPAddr("udp",":42069")
	if err != nil {
		log.Fatalf("Error resolving UDP addr: %s", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Error dialing UDP: %s", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading line: %s", err)
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatalf("Error writing line to UDP: %s", err)
		}
	}
}