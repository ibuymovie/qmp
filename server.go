package qmp

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"qmp/Handshake"
	"qmp/Message"
)

var ServerVersion = 1

type Server struct {
	port       int
	connectors []*Connector
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			return err
		}

		go func() {
			err := s.createConnector(conn)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
}

func (s *Server) createConnector(conn net.Conn) error {
	err := s.handshakeWithClient(bufio.NewWriter(conn), bufio.NewReader(conn))
	if err != nil {
		return err
	}

	c := NewConnector(conn)

	go c.RunRead()
	go c.RunWrite()

	s.connectors = append(s.connectors, c)

	return nil
}

func (s *Server) handshakeWithClient(w *bufio.Writer, r *bufio.Reader) error {
	// Recv C
	c, err := Handshake.DecodeSC(r)
	if err != nil {
		return err
	}

	if c.Version > byte(ServerVersion) {
		return errors.New("version non supported")
	}

	// Send S
	s0 := Handshake.NewSC(byte(ServerVersion))

	if err := s0.EncodeSC(w); err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) SendMessageToAll(message *Message.Message) {
	for _, connector := range s.connectors {
		connector.Write <- message
	}
}
