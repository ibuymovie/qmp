package qmp

import (
	"bufio"
	"fmt"
	"net"
	"qmp/Message"
)

type Connector struct {
	con    net.Conn
	writer *bufio.Writer
	reader *bufio.Reader
	write  chan *Message.Message
	read   chan *Message.Message
}

func NewConnector(con net.Conn) *Connector {
	return &Connector{
		con:    con,
		writer: bufio.NewWriter(con),
		reader: bufio.NewReader(con),
		write:  make(chan *Message.Message),
		read:   make(chan *Message.Message),
	}
}

func (c *Connector) RunRead() {
	for {
		chunk, err := Message.DecodeMessage(c.reader)
		if err != nil {
			continue
		}

		fmt.Println(chunk)
	}
}

func (c *Connector) RunWrite() {
	for chunk := range c.write {
		err := chunk.EncodeMessage(c.writer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = c.writer.Flush()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
