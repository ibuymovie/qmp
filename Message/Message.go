package Message

import (
	"io"
	"qmp/utils"
)

type Message struct {
	Type         uint16 // 2 bytes //
	BodyLength   uint32 // 4 bytes //
	HeaderLength uint32 // 4 bytes //
	Header       *Header
	Body         []byte
}

func DecodeMessage(r io.Reader) (*Message, error) {
	m := &Message{}

	buf := make([]byte, 8)

	if _, err := io.ReadAtLeast(r, buf[:2], 2); err != nil {
		return nil, err
	}
	m.Type = uint16(utils.ByteArrayToUint32(buf[:2]))

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

func (bh *Message) EncodeMessage(writer io.Writer) error {
	buf := make([]byte, 8)

	buf[0] = byte(bh.Type >> 8)
	buf[1] = byte(bh.Type & 0xff)

	_, err := writer.Write(buf[:2])
	if err != nil {
		return err
	}

	copy(buf[0:4], utils.Uint32ToByteArray(bh.BodyLength, 4))

	_, err = writer.Write(buf[:4])
	if err != nil {
		return err
	}

	copy(buf[0:4], utils.Uint32ToByteArray(bh.HeaderLength, 4))

	_, err = writer.Write(buf[:4])
	if err != nil {
		return err
	}

	return nil
}
