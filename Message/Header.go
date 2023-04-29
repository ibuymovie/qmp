package Message

import (
	"github.com/StounhandJ/go-amf"
	"io"
)

type Header struct {
	bodyLength int
	body       []byte
	Data       interface{}
}

func DecodeHeader(r io.Reader, headerLength uint32) (*Header, error) {
	if headerLength == 0 {
		return nil, nil
	}
	h := &Header{}

	buf := make([]byte, headerLength)

	if _, err := io.ReadAtLeast(r, buf[:], int(headerLength)); err != nil {
		return nil, err
	}

	h.body = buf[:]

	result, _, err := amf.DecodeAMF0(h.body)

	if err != nil {
		return nil, err
	}

	h.Data = result

	return h, nil
}

func (h *Header) EncodeHeader(w io.Writer) error {
	_, err := amf.EncodeAMF0(w, h.Data)

	return err
}
