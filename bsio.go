package bsio

import (
	//	"errors"
	"io"

	"fmt"
	"github.com/pkg/errors"
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
	if o == 0 && l == 8 {
		return b[0], nil
	}
	return 0, nil
}

type bigEndian struct{}

type Reader struct {
	buf   []byte
	rd    io.Reader
	order ByteOrder
	rp    int  // read pointer
	last  byte // last 1-byte data
}

func NewReader(r io.Reader, order ByteOrder) *Reader {
	b := new(Reader)
	b.rd = r
	b.order = order
	return b
}

// order: bsdecoder.LittleEndian
// data: pointer to data that will be output
// length: specify read length in bit
func (this *Reader) Read(data interface{}, length int) error {
	s := maxDataSize(data)
	if s == 0 {
		return errors.New("given type is not supported.")
	}
	if s < length {
		return errors.New(fmt.Sprintf("Given length (0x%x) is longer than given data type.", length))
	}

	fmt.Printf("size of data: %d bit(s).\n", s)

	if s%8 == 0 {
		// fast path: just read by io.ReadFull
		var b [8]byte
		bs := b[:(s / 8)]
		if _, err := io.ReadFull(this.rd, bs); err != nil {
			return errors.Wrap(err, "Failed to read by io.ReadFull")
		}

		switch data := data.(type) {
		case *uint8:
			var err error
			*data, err = this.order.Uint8(bs, 0, 8)
			if err != nil {
				return errors.Wrap(err, "Failed to convert to uint8")
			}
		default:
			return errors.New("Other than uint8 is not supported for now.")
		}

		return nil

	} else {
		return this.read(data, length)
	}

	return nil
}

func (this *Reader) read(data interface{}, length int) error {
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
