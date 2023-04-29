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

func TestDecodeEmptyHeader(t *testing.T) {
	var b bytes.Buffer

	reader := bufio.NewReader(&b)

	headers := map[string]interface{}{}

	header, headerData, err := Message.DecodeHeader(reader, uint32(0))

	assert.Equal(t, err, nil)
	assert.Equal(t, len(header), 0)
	assert.Equal(t, headerData, headers)
}

func TestEncodeHeader(t *testing.T) {

	header := map[string]interface{}{
		"App": "live",
		"Vid": "1",
	}

	_, err := Message.EncodeHeader(header)

	assert.Equal(t, err, nil)
}
