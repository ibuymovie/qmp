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
	Write  chan *Message.Message
	Read   chan *Message.Message
}

func NewConnector(con net.Conn) *Connector {
	return &Connector{
		con:    con,
		writer: bufio.NewWriter(con),
		reader: bufio.NewReader(con),
		Write:  make(chan *Message.Message),
		Read:   make(chan *Message.Message),
	}
}

func (c *Connector) RunRead() {
	for {
		message, err := Message.DecodeMessage(c.reader)
		if err != nil {
			continue
		}

		fmt.Println(message)
	}
}

func (c *Connector) RunWrite() {
	for chunk := range c.Write {
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
