package main

import (
	"fmt"
	"io/ioutil"
)

type bitStream struct {
	nibbles []uint8
	index   int
	end     []int
	err     error
}

func (b *bitStream) PushTruncation(n int) {
	e := b.index + n
}

func readDay16() (*bitStream, error) {
	f, err := ioutil.ReadFile("day16.txt")
	if err != nil {
		return nil, err
	}
	var nibbles []uint8
	for _, b := range f {
		if b >= '0' && b <= '9' {
			nibbles = append(nibbles, uint8(b-'0'))
		} else if b >= 'A' && b <= 'F' {
			nibbles = append(nibbles, uint8(10+b-'A'))
		} else if b != '\n' {
			return nil, fmt.Errorf("bad hex %c", b)
		}
	}
	return &bitStream{
		nibbles: nibbles,
		index:   -1,
		err:     nil,
	}, nil
}

// EOF reports whether next Scan() will fail
// because the stream is over.
func (b *bitStream) EOF() bool {
	return b.index+1 >= 4*len(b.nibbles)
}

func (b *bitStream) Scan() bool {
	if b.err != nil {
		return false
	}
	b.index++
	if b.index < 0 || b.index >= 4*len(b.nibbles) {
		return false
	}
	return true
}
func (b *bitStream) Err() error {
	return b.err
}

func (b *bitStream) Bit() bool {
	if b.err != nil || b.index < 0 || b.index >= 4*len(b.nibbles) {
		return false
	}
	n := b.nibbles[b.index/4]
	n >>= 3 - (b.index % 4)
	return (n & 1) == 1
}

func (b *bitStream) SetError(err error) {
	if b.err == nil {
		b.err = err
	}
}

func (b *bitStream) ScanBit() bool {
	if b.EOF() {
		b.SetError(fmt.Errorf("expected bit but found EOF"))
	}
	b.Scan()
	return b.Bit()
}

func readFixedInt(b *bitStream, n int) uint64 {
	var r uint64
	for i := 0; i < n; i++ {
		bit := b2i(b.ScanBit())
		r |= uint64(bit) << (n - i)
	}
	return r
}

type packet struct {
	version, id uint64
	literal     uint64
	packets     []packet
}

func readLiteral(b *bitStream) (uint64, error) {
	var r uint64
	nibs := 0
	for {
		prefix := b.ScanBit()
		nibble := readFixedInt(b, 4)
		nibs++
		if nibs > 16 {
			return 0, fmt.Errorf("literal doesn't fit in uint64")
		}
		r = (r << 4) | nibble
		if !prefix {
			return r, b.Err()
		}
	}
}

func readPacket(b *bitStream) (*packet, error) {
	if b.EOF() {
		return nil, nil
	}
	ver := readFixedInt(b, 3)
	id := readFixedInt(b, 3)
	if b.Err() != nil {
		return nil, b.Err()
	}
	switch id {
	case 4:
		lit, err := readLiteral(b)
		if err != nil {
			return nil, err
		}
		return &packet{
			version: ver,
			id:      id,
			literal: lit,
		}, nil
	default: // an operator
		lengthID := b.ScanBit()
		if lengthID {
			bitLen := readFixedInt(b, 15)
			b.PushTruncation(bitLen)
			defer b.PopTruncation()
			packets, err := readPackets(b, -1)
			if err != nil {
				return nil, err
			}
			return &packet{
				version: ver,
				id:      id,
				packets: packets,
			}, b.Err()
		}
	}
}

func day16() error {
	stream, err := readDay16()
	if err != nil {
		return err
	}
	for {
		p, err := readPacket(stream)
		if err != nil {
			return err
		}
		if p == nil {
			break
		}
		fmt.Println(p)
	}
	return stream.Err()
}

func init() {
	RegisterDay(16, day16)
}
