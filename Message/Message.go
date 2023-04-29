package Message

import (
	"io"
	"qmp/utils"
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
	Header       *Header
	body         []byte
	BodyData     interface{}
}

func DecodeMessage(r io.Reader) (*Message, error) {
	m, err := DecodeSetup(r)
	if err != nil {
		return nil, err
	}

	m.Header, err = DecodeHeader(r, m.HeaderLength)
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

	if err := mes.EncodeSetup(writer); err != nil {
		return err
	}

	if err := mes.Header.EncodeHeader(writer); err != nil {
		return err
	}

	if err := EncodeBody(writer, mes.BodyData, mes.MessageType); err != nil {
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
