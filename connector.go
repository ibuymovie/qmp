package qmp

import (
	"bufio"
	"fmt"
	"github.com/ibuymovie/qmp/Message"
	"net"
)

type Connector struct {
	con    net.Conn
	writer *bufio.Writer
	reader *bufio.Reader
	Write  chan *Message.Message
	Read   chan *Message.Message
	Close  chan bool
}

func newConnector(con net.Conn) *Connector {
	return &Connector{
		con:    con,
		writer: bufio.NewWriter(con),
		reader: bufio.NewReader(con),
		Write:  make(chan *Message.Message, 100),
		Read:   make(chan *Message.Message, 100),
		Close:  make(chan bool, 1),
	}
}

func (c *Connector) RunRead() {
	for {
		message, err := Message.DecodeMessage(c.reader)
		if err != nil {
			continue
		}

		c.Read <- message
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
			c.CloseConn()
			return
		}
	}
}

func (c *Connector) CloseConn() {
	_ = c.con.Close()
	close(c.Write)
	close(c.Read)
	c.Close <- true
}
