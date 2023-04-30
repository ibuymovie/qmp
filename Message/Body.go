package Message

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/StounhandJ/go-amf"
	"io"
)

func DecodeBody(r io.Reader, bodyLength uint32, messageType MessageType) ([]byte, interface{}, error) {

	buf := make([]byte, bodyLength)

	if _, err := io.ReadAtLeast(r, buf[:], int(bodyLength)); err != nil {
		return nil, nil, err
	}

	switch messageType {
	case Empty:
		return buf, nil, nil
	case String:
		return buf, string(buf), nil
	case Amf0:
		data, _, err := amf.DecodeAMF0(buf)
		if err != nil {
			return nil, nil, err
		}
		return buf, data, nil
	case Byte:
		return buf, buf, nil
	default:
		return nil, nil, errors.New("message type not allowed")
	}
}

func EncodeBody(data interface{}, messageType MessageType) ([]byte, error) {
	switch messageType {
	case Empty:
		return []byte{}, nil
	case String:
		return []byte(data.(string)), nil
	case Amf0:
		var b bytes.Buffer

		writer := bufio.NewWriter(&b)
		reader := bufio.NewReader(&b)
		length, err := amf.EncodeAMF0(writer, data)
		if err != nil {
			return nil, err
		}
		_ = writer.Flush()

		buf := make([]byte, length)
		_, _ = reader.Read(buf)
		return buf, nil
	case Byte:
		return data.([]byte), nil
	default:
		return nil, errors.New("message type not allowed")
	}
}
