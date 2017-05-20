package bsio

import (
	//	"errors"
	"io"

	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
)

type ByteOrder interface {
	Uint8([]byte, uint, uint) (uint8, error)
	// Uint16([]byte, uint, uint) (uint16, error)
	Uint32([]byte, uint, uint) (uint32, error)
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

func (littleEndian) Uint32(b []byte, o uint, l uint) (uint32, error) {
	if o == 0 {
		switch l {
		case 8:
			return uint32(b[0]), nil
		case 16:
			return uint32(b[0]) | (uint32(b[1]) << 8), nil
		case 24:
			return uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16), nil
		case 32:
			return uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24), nil
		}
	}
	return 0, nil
}

type bigEndian struct{}

type Reader struct {
	buf   []byte
	rd    io.Reader
	order ByteOrder

	last byte // last 1-byte data
	rp   uint // read pointer of last data. 0-7
}

func NewReader(r io.Reader, order ByteOrder) *Reader {
	b := new(Reader)
	b.rd = r
	b.order = order
	b.rp = 0
	return b
}

// order: bsdecoder.LittleEndian
// data: pointer to data that will be output
// length: specify read length in bit
func (this *Reader) Read(data interface{}, length uint) error {
	s := maxDataSize(data)
	if s == 0 {
		return errors.New("given type is not supported.")
	}
	if uint(s) < length {
		return errors.New(fmt.Sprintf("Given length (0x%x) is longer than given data type.", length))
	}

	fmt.Printf("size of data: %d bit(s).\n", s)

	if s%8 == 0 && this.rp == 0 && length%8 == 0 {
		return this.readBytes(data, length)
	} else {
		return this.readBit(data, length)
	}

	return nil
}

func (this *Reader) readBytes(data interface{}, length uint) error {
	var b [8]byte
	bs := b[:(length / 8)]
	if _, err := io.ReadFull(this.rd, bs); err != nil {
		return errors.Wrap(err, "Failed to read by io.ReadFull")
	}

	switch data := data.(type) {
	case *uint8:
		var err error
		*data, err = this.order.Uint8(bs, 0, length)
		if err != nil {
			return errors.Wrap(err, "Failed to convert to uint8")
		}
	case *uint32:
		var err error
		*data, err = this.order.Uint32(bs, 0, length)
		if err != nil {
			return errors.Wrap(err, "Failed to convert to uint8")
		}
	default:
		return errors.New("Other than uint8 is not supported for now.")
	}

	return nil
}

// const (
// 	l0 = 0x00
// 	l1 = 0x80
// 	l2 = 0xC0
// 	l3 = 0xE0
// 	l4 = 0xF0
// 	l5 = 0xF8
// 	l6 = 0xFC
// 	l7 = 0xFE
// 	l8 = 0xFF
// )

// var lemask = []byte{l0, l1, l2, l3, l4, l5, l6, l7, l8}

func (this *Reader) readBit(data interface{}, length uint) error {

	var sum int64
	for {
		if this.rp == 0 {
			// no data to process: load 1 byte
			err := binary.Read(this.rd, binary.LittleEndian, &(this.last))
			if err != nil {
				return errors.Wrap(err, "Failed to read next byte")
			}
			// fmt.Printf("new last: 0x%x\n", this.last)
		}

		// read length
		r_length := length
		remain := uint(8) - this.rp
		if length > remain {
			r_length = remain
		}

		// note: littleEndian
		d := int64((this.last << this.rp) >> (8 - r_length))
		sum = (sum << r_length) + d

		fmt.Printf("last: %02x rp: %x r_len: %x d: %x sum: %x\n", this.last, this.rp, r_length, d, sum)

		// increment read pointer
		this.rp += r_length
		length -= r_length

		if this.rp == 8 {
			this.rp = 0
			this.last = 0
		}

		if length == 0 {
			break
		}

	}

	// store data to each types
	switch data := data.(type) {
	case *uint8:
		fmt.Printf("read: 0x%x\n", sum)
		*data = uint8(sum)
	case *uint32:
		*data = uint32(sum)
	default:
		return errors.New("Other than uint8 is not supported for now.")
	}

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
