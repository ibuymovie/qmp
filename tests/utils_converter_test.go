package tests

import (
	"github.com/go-playground/assert/v2"
	"qmp/utils"
	"testing"
)

func TestUintToBinaryArray(t *testing.T) {
	type args struct {
		dataNumber uint32
		dataLength int
		expected   []byte
	}

	tests := []args{
		{
			1,
			2,
			[]byte{0, 1},
		},
		{
			340,
			2,
			[]byte{1, 84},
		},
		{
			340,
			4,
			[]byte{0, 0, 1, 84},
		},
		{
			34690,
			3,
			[]byte{0, 135, 130},
		},
		{
			563034,
			3,
			[]byte{8, 151, 90},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, utils.Uint32ToByteArray(tt.dataNumber, tt.dataLength), tt.expected)
	}
}
