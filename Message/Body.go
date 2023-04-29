package Message

import (
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
