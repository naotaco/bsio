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

type bsData struct {
	seq []byte   // raw byte sequence
	b   [][]byte // expected data for 1~7 bit reading
}

var bsTestData = []bsData{
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
	for _, d_ := range bsTestData {
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
