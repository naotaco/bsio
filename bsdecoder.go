package bsdecoder

import (
	"errors"
	"io"

	"fmt"
)

type ByteOrder interface {
	Uint8([]byte, uint, uint) (uint8, error)
	// Uint16([]byte, uint, uint) (uint16, error)
	// Uint32([]byte, uint, uint) (uint32, error)
	// Uint64([]byte, uint, uint) (uint64, error)
	// PutUint8([]byte, uint8)
	// PutUint16([]byte, uint16)
	// PutUint32([]byte, uint32)
	// PutUint64([]byte, uint64)
	// String() string
}

var LittleEndian littleEndian
var BigEndian bigEndian

type littleEndian struct{}

func (littleEndian) Uint8(b []byte, o uint, l uint) (uint8, error) {

	return 0, nil
}

type bigEndian struct{}

// order: bsdecoder.LittleEndian
// data: pointer to data that will be output
// length: specify read length in bit
func Read(r io.Reader, order ByteOrder, data interface{}, length int) error {
	s := maxDataSize(data)
	if s == 0 {
		return errors.New("given type is not supported.")
	}
	if s < length {
		return errors.New(fmt.Sprintf("Given length (0x%x) is longer than given data type.", length))
	}

	fmt.Printf("size of data: %d bit(s).\n", s)
	return nil
}

// returns size in BIT.
func maxDataSize(data interface{}) int {
	switch data := data.(type) {
	case bool, int8, uint8, *bool, *int8, *uint8:
		return 8
	case []int8:
		return 8 * len(data)
	case []uint8:
		return 8 * len(data)
	case int16, uint16, *int16, *uint16:
		return 16
	case []int16:
		return 16 * len(data)
	case []uint16:
		return 16 * len(data)
	case int32, uint32, *int32, *uint32:
		return 32
	case []int32:
		return 32 * len(data)
	case []uint32:
		return 32 * len(data)
	case int64, uint64, *int64, *uint64:
		return 64
	case []int64:
		return 64 * len(data)
	case []uint64:
		return 64 * len(data)
	}
	return 0
}
