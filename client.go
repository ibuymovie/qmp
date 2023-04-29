package qmp

import (
	"bufio"
	"errors"
	"net"
	"qmp/Handshake"
	"qmp/Message"
)

var ClientVersion = 1

type Client struct {
	address   string
	connector *Connector
	Messages  chan *Message.Message
}

func NewClient(address string) *Client {
	return &Client{
		address:  address,
		Messages: make(chan *Message.Message),
	}
}

func (c *Client) Run() error {
	ip, err := net.Dial("tcp", c.address)
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
