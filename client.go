package qmp

import (
	"bufio"
	"errors"
	"github.com/ibuymovie/qmp/Handshake"
	"github.com/ibuymovie/qmp/Message"
	"net"
)

var ClientVersion = 1

type Client struct {
	Address   string
	connector *Connector
	Messages  chan *Message.Message
}

func NewClient(address string) *Client {
	return &Client{
		Address:  address,
		Messages: make(chan *Message.Message),
	}
}

func (c *Client) Run() error {
	ip, err := net.Dial("tcp", c.Address)
	if err != nil {
		return err
	}

	err = c.handshakeWithSerer(bufio.NewWriter(ip), bufio.NewReader(ip))
	if err != nil {
		return err
	}

	c.connector = newConnector(ip)

	go c.connector.RunRead()
	go c.connector.RunWrite()
	go func() {
		for {
			c.Messages <- <-c.connector.Read
		}
	}()

	return nil
}

func (c *Client) handshakeWithSerer(w *bufio.Writer, r *bufio.Reader) error {
	// Send C
	c0 := Handshake.NewSC(byte(ClientVersion))

	if err := c0.EncodeSC(w); err != nil {
		return err
	}
	err := w.Flush()

	if err != nil {
		return err
	}

	// Recv S
	s0, err := Handshake.DecodeSC(r)
	if err != nil {
		return err
	}

	if s0.Version > byte(ClientVersion) {
		return errors.New("version non supported")
	}

	return nil
}

func (c *Client) Close() {
	c.connector.CloseConn()
}

func (c *Client) SendMessageToServer(message *Message.Message) {
	go func() {
		c.connector.Write <- message
	}()
}
