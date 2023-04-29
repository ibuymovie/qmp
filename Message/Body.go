package Message

import (
	"bufio"
	"bytes"
	"encoding/json"
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
	case Json:
		data := make(map[string]interface{})
		err := json.Unmarshal(buf, &data)
		if err != nil {
			return nil, nil, err
		}
		return buf, data, nil
	case Amf0:
		data, _, err := amf.DecodeAMF0(buf)
		if err != nil {
			return nil, nil, err
		}
		return buf, data, nil
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
	case Json:
		body, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		return body, nil
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
	default:
		return nil, errors.New("message type not allowed")
	}
}
