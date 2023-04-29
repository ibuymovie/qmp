package tests

import (
	"bufio"
	"bytes"
	"github.com/StounhandJ/go-amf"
	"github.com/go-playground/assert/v2"
	"qmp/Message"
	"testing"
)

func TestDecodeHeader(t *testing.T) {
	var b bytes.Buffer

	writer := bufio.NewWriter(&b)
	reader := bufio.NewReader(&b)

	headers := map[string]interface{}{
		"App": "live",
		"Vid": "1",
	}

	length, _ := amf.EncodeAMF0(writer, headers)
	_ = writer.Flush()

	header, err := Message.DecodeHeader(reader, uint32(length))

	assert.Equal(t, err, nil)
	assert.Equal(t, header.Data, headers)
}

func TestDecodeNumberHeader(t *testing.T) {
	var b bytes.Buffer

	writer := bufio.NewWriter(&b)
	reader := bufio.NewReader(&b)

	writer.Write([]byte{0, 64, 106, 192, 0, 0, 0, 0, 0})
	_ = writer.Flush()

	header, err := Message.DecodeHeader(reader, 9)

	assert.Equal(t, err, nil)
	assert.Equal(t, header.Data, float64(214))
}

func TestEncodeHeader(t *testing.T) {
	var b bytes.Buffer

	writer := bufio.NewWriter(&b)
	reader := bufio.NewReader(&b)

	header := Message.Header{Data: 214}
	err := header.EncodeHeader(writer)
	_ = writer.Flush()

	res, _ := reader.ReadBytes(1)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []byte{0, 64, 106, 192, 0, 0, 0, 0, 0})
}
