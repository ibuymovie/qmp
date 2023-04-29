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

func EncodeBody(w io.Writer, data interface{}, messageType MessageType) error {
	switch messageType {
	case Empty:
		return nil
	case String:
		_, _ = w.Write([]byte(data.(string)))
		return nil
	case Json:
		body, err := json.Marshal(data)
		if err != nil {
			return err
		}
		_, _ = w.Write(body)
		return nil
	case Amf0:
		_, err := amf.EncodeAMF0(w, data)
		return err
	default:
		return errors.New("message type not allowed")
	}
}
