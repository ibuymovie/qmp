package Handshake

import "io"

type SC struct {
	Version byte
}

func NewSC(version byte) *SC {
	return &SC{
		Version: version,
	}
}

func DecodeSC(reader io.Reader) (*SC, error) {
	sc := &SC{}

	buf := [1]byte{}

	if _, err := io.ReadAtLeast(reader, buf[:], 1); err != nil {
		return nil, err
	}
	sc.Version = buf[0]

	return sc, nil
}

func (sc *SC) EncodeSC(writer io.Writer) error {
	buf := [1]byte{sc.Version}

	_, err := writer.Write(buf[:])
	if err != nil {
		return err
	}

	return nil
}
