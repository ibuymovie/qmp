package tests

import (
	"bufio"
	"bytes"
	"github.com/StounhandJ/go-amf"
	"github.com/go-playground/assert/v2"
	"github.com/ibuymovie/qmp/Message"
	"testing"
)

func TestDecodeBodyAmf0(t *testing.T) {
	var b bytes.Buffer

	writer := bufio.NewWriter(&b)
	reader := bufio.NewReader(&b)

	sendData := map[string]interface{}{
		"App": "live",
		"Vid": "1",
	}

	length, _ := amf.EncodeAMF0(writer, sendData)
	_ = writer.Flush()

	body, data, err := Message.DecodeBody(reader, uint32(length), Message.Amf0)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(body), length)
	assert.Equal(t, data, sendData)
}

func TestDecodeBodyString(t *testing.T) {
	var b bytes.Buffer

	writer := bufio.NewWriter(&b)
	reader := bufio.NewReader(&b)

	message := "message"
	length, _ := writer.Write([]byte(message))
	_ = writer.Flush()

	body, data, err := Message.DecodeBody(reader, uint32(length), Message.String)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(body), length)
	assert.Equal(t, data, message)
}

func TestDecodeBodyEmpty(t *testing.T) {
	var b bytes.Buffer

	reader := bufio.NewReader(&b)

	body, data, err := Message.DecodeBody(reader, uint32(0), Message.Empty)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(body), 0)
	assert.Equal(t, data, nil)
}
