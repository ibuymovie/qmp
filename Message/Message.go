package Message

import (
	"github.com/ibuymovie/qmp/utils"
	"io"
)

type MessageType uint16

var (
	Empty  MessageType = 0
	String MessageType = 1
	Json   MessageType = 2
	Amf0   MessageType = 3
)

type Message struct {
	MessageType  MessageType // 2 bytes //
	BodyLength   uint32      // 4 bytes //
	HeaderLength uint32      // 4 bytes //
	header       []byte
	HeaderData   interface{}
	body         []byte
	BodyData     interface{}
}

func NewMessage(messageType MessageType, BodyData interface{}) *Message {
	return &Message{
		MessageType: messageType,
		BodyData:    BodyData,
	}
}

func DecodeMessage(r io.Reader) (*Message, error) {
	m, err := DecodeSetup(r)
	if err != nil {
		return nil, err
	}

	m.header, m.HeaderData, err = DecodeHeader(r, m.HeaderLength)
	if err != nil {
		return nil, err
	}

	m.body, m.BodyData, err = DecodeBody(r, m.BodyLength, m.MessageType)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func DecodeSetup(r io.Reader) (*Message, error) {
	m := &Message{}

	buf := make([]byte, 8)

	if _, err := io.ReadAtLeast(r, buf[:2], 2); err != nil {
		return nil, err
	}
	m.MessageType = MessageType(utils.ByteArrayToUint32(buf[:2]))

	if _, err := io.ReadAtLeast(r, buf[:4], 4); err != nil {
		return nil, err
	}
	m.BodyLength = utils.ByteArrayToUint32(buf[:4])

	if _, err := io.ReadAtLeast(r, buf[:4], 4); err != nil {
		return nil, err
	}
	m.HeaderLength = utils.ByteArrayToUint32(buf[:4])

	return m, nil
}

func (mes *Message) EncodeMessage(writer io.Writer) error {
	var header, body []byte
	var err error

	if header, err = EncodeHeader(mes.HeaderData); err != nil {
		return err
	}
	mes.header = header
	mes.HeaderLength = uint32(len(header))

	if body, err = EncodeBody(mes.BodyData, mes.MessageType); err != nil {
		return err
	}
	mes.body = body
	mes.BodyLength = uint32(len(body))

	if err := mes.EncodeSetup(writer); err != nil {
		return err
	}

	if _, err := writer.Write(mes.header); err != nil {
		return err
	}

	if _, err := writer.Write(mes.body); err != nil {
		return err
	}

	return nil
}

func (mes *Message) EncodeSetup(writer io.Writer) error {
	buf := make([]byte, 8)

	buf[0] = byte(mes.MessageType >> 8)
	buf[1] = byte(mes.MessageType & 0xff)

	_, err := writer.Write(buf[:2])
	if err != nil {
		return err
	}

	copy(buf[0:4], utils.Uint32ToByteArray(mes.BodyLength, 4))

	_, err = writer.Write(buf[:4])
	if err != nil {
		return err
	}

	copy(buf[0:4], utils.Uint32ToByteArray(mes.HeaderLength, 4))

	_, err = writer.Write(buf[:4])
	if err != nil {
		return err
	}

	return nil
}
