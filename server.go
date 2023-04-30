package qmp

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ibuymovie/qmp/Handshake"
	"github.com/ibuymovie/qmp/Message"
	"net"
)

var ServerVersion = 1

type Server struct {
	port       int
	connectors []*Connector
	listener   net.Listener
	Messages   chan *Message.Message
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Run() error {
	var err error

	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			return err
		}
		fmt.Println("New connection", conn.RemoteAddr())

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

	c := newConnector(conn)

	go c.RunRead()
	go c.RunWrite()
	go func() {
		<-c.Close
		s.closeConnector(c)
	}()
	go func() {
		for {
			s.Messages <- <-c.Read
		}
	}()

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

func (s *Server) closeConnector(connector *Connector) {
	for i, c := range s.connectors {
		if c == connector {
			s.connectors = append(s.connectors[:i], s.connectors[i+1:]...)
			fmt.Println("Close connection", connector.con.RemoteAddr())
		}
	}
}

func (s *Server) SendMessageToAll(message *Message.Message) {
	go func() {
		for _, connector := range s.connectors {
			connector.Write <- message
		}
	}()
}

func (s *Server) Close() {
	for _, connector := range s.connectors {
		connector.CloseConn()
	}

	_ = s.listener.Close()
}
