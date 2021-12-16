package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type bitStream struct {
	nibbles []uint8
	index   int
	end     []int
	err     error
}

func (b *bitStream) PushTruncation(n int) {
	e := b.index + n + 1
	b.end = append(b.end, e)
}

func (b *bitStream) PopTruncation() {
	b.end = b.end[:len(b.end)-1]
}

func newBitStream(f []byte) (*bitStream, error) {
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
		end:     []int{len(nibbles) * 4},
		err:     nil,
	}, nil
}

func readDay16() (*bitStream, error) {
	f, err := ioutil.ReadFile("day16.txt")
	if err != nil {
		return nil, err
	}
	return newBitStream(f)
}

func (b *bitStream) eofPos() int {
	return b.end[len(b.end)-1]
}

// EOF reports whether next Scan() will fail
// because the stream is over.
func (b *bitStream) EOF() bool {
	return b.index+1 >= b.eofPos()
}

func (b *bitStream) Scan() bool {
	if b.err != nil {
		return false
	}
	if b.index >= b.eofPos() {
		return false
	}
	b.index++
	if b.index < 0 || b.index >= b.eofPos() {
		return false
	}
	return true
}
func (b *bitStream) Err() error {
	return b.err
}

func (b *bitStream) Bit() bool {
	if b.err != nil || b.index < 0 || b.index >= b.eofPos() {
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
		r |= uint64(bit) << (n - i - 1)
	}
	return r
}

type packet struct {
	version, id uint64
	literal     uint64
	packets     []*packet
}

func (p packet) String() string {
	if p.id == 4 {
		return fmt.Sprintf("[%d:%d lit=%d]", p.version, p.id, p.literal)
	}
	var subpackets []string
	for _, p := range p.packets {
		subpackets = append(subpackets, p.String())
	}
	return fmt.Sprintf("[%d:%d %s]", p.version, p.id, strings.Join(subpackets, " "))
}

func readLiteral(b *bitStream) (uint64, error) {
	var r uint64
	nibs := 0
	for {
		prefix := b.ScanBit()
		if b.Err() != nil {
			return 0, fmt.Errorf("literal prefix bit read after %d nibbles: %v", nibs, b.Err())
		}
		nibble := readFixedInt(b, 4)
		if b.Err() != nil {
			return 0, fmt.Errorf("nibble read after %d nibbles: %v", nibs, b.Err())
		}
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

func readPackets(b *bitStream, n int) ([]*packet, error) {
	var r []*packet
	for i := 0; n == -1 || i < n; i++ {
		pack, err := readPacket(b)
		if err != nil {
			return nil, err
		}
		if pack == nil {
			if n != -1 {
				return nil, fmt.Errorf("expected %d packets, but found only %d", n, len(r))
			}
			return r, b.Err()
		}
		r = append(r, pack)
	}
	return r, b.Err()
}

func packetVersionSum(p *packet) int {
	s := int(p.version)
	for _, c := range p.packets {
		s += packetVersionSum(c)
	}
	return s
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
		if !lengthID {
			bitLen := int(readFixedInt(b, 15))
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
		} else {
			packLen := int(readFixedInt(b, 11))
			packets, err := readPackets(b, packLen)
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

func packetEvaluate(p *packet) int64 {
	if p.id == 0 {
		var s int64
		for _, x := range p.packets {
			s += packetEvaluate(x)
		}
		return s
	} else if p.id == 1 {
		var s int64 = 1
		for _, x := range p.packets {
			s *= packetEvaluate(x)
		}
		return s
	} else if p.id == 2 {
		var s int64 = 1 << 62
		for _, x := range p.packets {
			i := packetEvaluate(x)
			if i < s {
				s = i
			}
		}
		return s
	} else if p.id == 3 {
		var s int64 = -(1 << 62)
		for _, x := range p.packets {
			i := packetEvaluate(x)
			if i > s {
				s = i
			}
		}
		return s
	} else if p.id == 4 {
		return int64(p.literal)
	} else if p.id == 5 || p.id == 6 || p.id == 7 {
		if len(p.packets) != 2 {
			panic(fmt.Sprintf("expected 2 packets in comparison but got %d", len(p.packets)))
		}
		left := packetEvaluate(p.packets[0])
		right := packetEvaluate(p.packets[1])
		if p.id == 5 {
			return int64(b2i(left > right))
		} else if p.id == 6 {
			return int64(b2i(left < right))
		} else if p.id == 7 {
			return int64(b2i(left == right))
		}
	}
	panic(fmt.Sprintf("unknown operator id %d", p.id))
}

func day16() error {
	stream, err := readDay16()
	if err != nil {
		return err
	}
	p, err := readPacket(stream)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("expected packet, but found none!?")
	}
	partPrint(1, packetVersionSum(p))
	partPrint(2, packetEvaluate(p))
	return stream.Err()
}

func init() {
	RegisterDay(16, day16)
}
