package Message

import (
	"bufio"
	"bytes"
	"github.com/StounhandJ/go-amf"
	"io"
)

func DecodeHeader(r io.Reader, headerLength uint32) ([]byte, map[string]interface{}, error) {
	if headerLength == 0 {
		return []byte{}, map[string]interface{}{}, nil
	}

	header := make([]byte, headerLength)

	if _, err := io.ReadAtLeast(r, header[:], int(headerLength)); err != nil {
		return nil, nil, err
	}

	headerData, _, err := amf.DecodeAMF0(header)

	if err != nil {
		return nil, nil, err
	}

	return header, headerData.(map[string]interface{}), nil
}

func EncodeHeader(header map[string]interface{}) ([]byte, error) {

	var b bytes.Buffer

	writer := bufio.NewWriter(&b)
	reader := bufio.NewReader(&b)
	length, err := amf.EncodeAMF0(writer, header)
	if err != nil {
		return nil, err
	}
	_ = writer.Flush()

	buf := make([]byte, length)
	_, _ = reader.Read(buf)
	return buf, nil
}
