package server

import (
	"log"
	"net"
)

type ServerState int

const (
	INITIALIZED ServerState = iota
	CLOSED
)

type Server struct {
	ServerState ServerState
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		return nil, err
	}
	server := &Server{
		ServerState: INITIALIZED,
		listener: ln,
	}
	server.listen()

	return server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.ServerState = CLOSED
	return nil
}

func (s *Server) listen() {
	for s.ServerState != CLOSED {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalf("could not open conn %s: %s\n", conn, err)
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	data := "HTTP/1.1 200 OK \r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"
	_, err := conn.Write([]byte(data))
	if err != nil {
		log.Fatalf("could not open conn %s: %s\n", conn, err)
	}
}