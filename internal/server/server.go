package server

import (
	"HTTPFROMTCP/internal/request"
	"HTTPFROMTCP/internal/response"
	"bytes"
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

func Serve(port int, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		return nil, err
	}
	server := &Server{
		ServerState: INITIALIZED,
		listener: ln,
	}
	server.listen(handler)

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

func (s *Server) listen(handler Handler) {
	for s.ServerState != CLOSED {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalf("could not open conn %s: %s\n", conn, err)
		}
		go s.handle(conn, handler)
	}
}

func (s *Server) handle(conn net.Conn, handler Handler) {
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatalf("error parsing request: %s\n", err)
	}

	var buf []byte
	buffer := bytes.NewBuffer(buf)
	handlerError := handler(buffer, req)

	err = response.WriteStatusLine(conn, response.StatusCode(handlerError.StatusCode))
	if err != nil {
		log.Fatalf("could not write Status Line: %s\n", err)
	}

	headers := response.GetDefaultHeaders(buffer.Len())
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Fatalf("could not write headers: %s", err)
	}

	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		log.Fatalf("could not write body: %s", err)
	}
}