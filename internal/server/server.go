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
	handler Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		return nil, err
	}
	server := &Server{
		ServerState: INITIALIZED,
		listener: ln,
		handler: handler,
	}
	go server.listen()

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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.BAD_REQ,
			Message:    err.Error(),
		}
		hErr.Write(conn)
	}

	buffer := bytes.NewBuffer([]byte{})
	handlerError := s.handler(buffer, req)
	if handlerError != nil {
		handlerError.Write(conn)
		return
	}

	response.WriteStatusLine(conn, response.OK)
	headers := response.GetDefaultHeaders(buffer.Len())
	response.WriteHeaders(conn, headers)

	conn.Write(buffer.Bytes())
}