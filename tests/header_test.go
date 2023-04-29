package tests

import (
	"bufio"
	"bytes"
	"github.com/StounhandJ/go-amf"
	"github.com/go-playground/assert/v2"
	"github.com/ibuymovie/qmp/Message"
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

	header, headerData, err := Message.DecodeHeader(reader, uint32(length))

	assert.Equal(t, err, nil)
	assert.Equal(t, len(header), length)
	assert.Equal(t, headerData, headers)
}

func TestDecodeNumberHeader(t *testing.T) {
	var b bytes.Buffer

	writer := bufio.NewWriter(&b)
	reader := bufio.NewReader(&b)

	length, _ := writer.Write([]byte{0, 64, 106, 192, 0, 0, 0, 0, 0})
	_ = writer.Flush()

	header, headerData, err := Message.DecodeHeader(reader, uint32(length))

	assert.Equal(t, err, nil)
	assert.Equal(t, len(header), length)
	assert.Equal(t, headerData, float64(214))
}

func TestEncodeHeader(t *testing.T) {

	header, err := Message.EncodeHeader(214)

	assert.Equal(t, err, nil)
	assert.Equal(t, header, []byte{0, 64, 106, 192, 0, 0, 0, 0, 0})
}
