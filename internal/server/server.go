package server

import (
	"HTTPFROMTCP/internal/response"
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
	err := response.WriteStatusLine(conn, response.OK)
	if err != nil {
		log.Fatalf("could not write Status Line: %s\n", err)
	}

	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		log.Fatalf("could not write headers: %s", err)
	}
}