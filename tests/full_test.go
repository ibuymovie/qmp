package tests

import (
	"bufio"
	"bytes"
	"github.com/go-playground/assert/v2"
	"qmp/Message"
	"testing"
)

func TestDecode(t *testing.T) {
	message := make([]byte, 10)

	// Message type
	message[1] = 0x01
	// Body length
	message[5] = 0x14
	// Header length
	message[9] = 0xa

	mes, err := Message.DecodeMessage(bytes.NewReader(message))

	assert.Equal(t, err, nil)
	assert.Equal(t, mes.Type, uint16(1))
	assert.Equal(t, mes.BodyLength, uint32(20))
	assert.Equal(t, mes.HeaderLength, uint32(10))
}

func TestDecodeBigNumber(t *testing.T) {
	message := make([]byte, 10)

	// Message type
	message[0] = 0x05
	// Body length
	message[4] = 0x21
	// Header length
	message[7] = 0x5
	message[9] = 0xa3

	mes, err := Message.DecodeMessage(bytes.NewReader(message))

	assert.Equal(t, err, nil)
	assert.Equal(t, mes.Type, uint16(1280))
	assert.Equal(t, mes.BodyLength, uint32(8448))
	assert.Equal(t, mes.HeaderLength, uint32(327843))
}

func TestEncode(t *testing.T) {
	message := Message.Message{
		Type:         1,
		BodyLength:   20,
		HeaderLength: 10,
	}

	var b bytes.Buffer

	writer := bufio.NewWriter(&b)
	reader := bufio.NewReader(&b)

	err := message.EncodeMessage(writer)
	writer.Flush()

	res, _ := reader.ReadBytes(10)

	assert.Equal(t, err, nil)
	assert.Equal(t, res[1], byte(1))
	assert.Equal(t, res[5], byte(20))
	assert.Equal(t, res[9], byte(10))
}
