package bsio

import (
	//	"bsio"
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"testing"
)

var testByteSeq = [][]byte{
	[]byte{0x12, 0xab, 0x22, 0xff, 0x00, 0x34, 0x09, 0x01, 0x90, 0xaa},
	[]byte{0},
	[]byte{},
}

type bsData8 struct {
	seq []byte   // raw byte sequence
	b   [][]byte // expected data for 1~7 bit reading
}

type bsData32 struct {
	seq []byte
	b   [][]uint32
}

var bsTestData8 = []bsData8{
	{
		//     0b00010010, 0b00110100
		[]byte{0x12, 0x34},
		[][]byte{
			[]byte{}, // dummy
			[]byte{0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0}, // 1-bit reading
			[]byte{0, 1, 0, 2, 0, 3, 1, 0},                         // 2-bit
			[]byte{0, 4, 4, 3, 2},                                  // ...
			[]byte{1, 2, 3, 4},
			[]byte{2, 8, 26},
			[]byte{4, 35},
			[]byte{9, 13},
		},
	},
}

var bsTestData32 = []bsData32{
	{
		//     0b01000010, 0b10010111, 0b01110110, 0b10000110
		[]byte{0x42, 0x97, 0x76, 0x86},
		[][]uint32{
			[]uint32{}, // dummy.
			[]uint32{}, // 1-bit reading
			[]uint32{}, // 2-bit
			[]uint32{}, // ...
			[]uint32{},
			[]uint32{},
			[]uint32{},
			[]uint32{},                       // 7
			[]uint32{0x42, 0x97, 0x76, 0x86}, // 8-bit
			[]uint32{0x85, 0x5d, 0x1b4},      // 9-bit
			[]uint32{0x10a, 0x177, 0x1a1},    // 10-bit
			[]uint32{0x214, 0x5dd},           // 11-bit
			[]uint32{},                       // 12
			[]uint32{},                       // 13
			[]uint32{},                       // 14
			[]uint32{0x214b, 0x5da1},         // 15
			[]uint32{0x4297, 0x7686},         // 16
			[]uint32{},                       // 17
		},
	},
}

func TestStdBinaryCompatibility(t *testing.T) {

	for _, b := range testByteSeq {
		b_ := b
		r1 := bytes.NewReader(b)
		r2 := bytes.NewReader(b_)
		bs := NewReader(r2, LittleEndian)
		for {
			var v1 byte
			var v2 byte
			err := binary.Read(r1, binary.LittleEndian, &v1)
			err2 := bs.Read(&v2, 8)
			// bsio.Read wraps io.EOF
			if err == io.EOF && errors.Cause(err2) == io.EOF {
				// succeed: both of Read functions should return EOF in same trial.
				break
			}
			if v1 != v2 {
				t.Logf("Unmatched read 8-bit data")
				t.Fail()
			} else {
				// OK.
				// fmt.Printf("v1 %x v2 %x\n", v1, v2)
			}
		}
	}
}

func TestUint8(t *testing.T) {
	for _, d_ := range bsTestData8 {
		for _, bit := range []uint{1, 2, 3, 4, 5, 6, 7} {

			d := d_
			r := bytes.NewReader(d.seq)
			bs := NewReader(r, LittleEndian)

			i := 0
			for {
				var v byte
				err := bs.Read(&v, bit)
				if errors.Cause(err) == io.EOF {
					if i != len(d.b[bit]) {
						t.Logf("Unmatched length: i: %d len: %d\n", i, len(d.b[bit]))
						t.Fail()
					}
					t.Logf("i: %d len: %d\n", i, len(d.b[bit]))
					break
				}

				if v != d.b[bit][i] {
					t.Logf("%d th data unmatched. act: %x expected: %x bit len: %x ", i, v, d.b[bit][i], bit)
					t.Fail()
				} else {
					// t.Logf("act: %x expected: %x\n", v, d.b1[i])
				}

				i++
			}
		}
	}
}

func TestUint32(t *testing.T) {
	for _, d_ := range bsTestData32 {
		for _, bit := range []uint{8, 9, 10, 11, 15} {

			d := d_
			r := bytes.NewReader(d.seq)
			bs := NewReader(r, LittleEndian)

			i := 0
			for {
				var v uint32
				err := bs.Read(&v, bit)
				if errors.Cause(err) == io.EOF {
					if i != len(d.b[bit]) {
						t.Logf("Unmatched length: i: %d len: %d\n", i, len(d.b[bit]))
						t.Fail()
					}
					t.Logf("i: %d len: %d\n", i, len(d.b[bit]))
					break
				}

				if v != d.b[bit][i] {
					t.Logf("%d th data unmatched. act: %x expected: %x bit len: %x ", i, v, d.b[bit][i], bit)
					t.Fail()
				} else {
					// t.Logf("act: %x expected: %x\n", v, d.b1[i])
				}

				i++
			}
		}
	}
}
