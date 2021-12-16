package main

import (
	"fmt"
	"reflect"
	"testing"
)

type spTC struct {
	in         string
	want       *packet
	wantVerSum int
	wantEV     int64
}

func TestStreamParser(t *testing.T) {
	cases := []spTC{
		{
			in:   "D2FE28",
			want: &packet{version: 6, id: 4, literal: 2021},
		},
		{
			in: "38006F45291200",
			want: &packet{version: 1, id: 6, packets: []*packet{
				{version: 6, id: 4, literal: 10},
				{version: 2, id: 4, literal: 20},
			}},
		},
		{
			in: "EE00D40C823060",
			want: &packet{version: 7, id: 3, packets: []*packet{
				{version: 2, id: 4, literal: 1},
				{version: 4, id: 4, literal: 2},
				{version: 1, id: 4, literal: 3},
			}},
		},
		{
			in: "8A004A801A8002F478",
			want: &packet{version: 4, id: 2, packets: []*packet{
				{version: 1, id: 2, packets: []*packet{
					{version: 5, id: 2, packets: []*packet{
						{version: 6, id: 4, literal: 15},
					}},
				}},
			}},
			wantVerSum: 16,
		},
		{
			in:         "620080001611562C8802118E34",
			wantVerSum: 12,
		},
		{
			in:         "C0015000016115A2E0802F182340",
			wantVerSum: 23,
		},
		{
			in:         "A0016C880162017C3686B18A3D4780",
			wantVerSum: 31,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("case[%d]=%q", i, c.in), func(t *testing.T) {
			b, err := newBitStream([]byte(c.in))
			if err != nil {
				t.Fatal(err)
			}
			got, err := readPacket(b)
			if err != nil {
				t.Fatal(err)
			}
			if c.want != nil {
				if !reflect.DeepEqual(got, c.want) {
					t.Errorf("readPacket() = %v, want %v", got, c.want)
				}
			}
			if c.wantVerSum != 0 {
				gotSum := packetVersionSum(got)
				if gotSum != c.wantVerSum {
					t.Errorf("packetVersionSum()=%d, want %d", gotSum, c.wantVerSum)
				}
			}
		})
	}

}

type evalTC struct {
	in     string
	wantEV int64
}

func TestPacketEval(t *testing.T) {
	cases := []evalTC{
		{
			in:     "C200B40A82",
			wantEV: 3,
		},
		{
			in:     "04005AC33890",
			wantEV: 54,
		},
		{
			in:     "880086C3E88112",
			wantEV: 7,
		},
		{
			in:     "CE00C43D881120",
			wantEV: 9,
		},
		{
			in:     "D8005AC2A8F0",
			wantEV: 1,
		},
		{
			in:     "F600BC2D8F",
			wantEV: 0,
		},
		{
			in:     "9C005AC2F8F0",
			wantEV: 0,
		},
		{
			in:     "9C0141080250320F1802104A08",
			wantEV: 1,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("case[%d]=%q", i, c.in), func(t *testing.T) {
			b, err := newBitStream([]byte(c.in))
			if err != nil {
				t.Fatal(err)
			}
			p, err := readPacket(b)
			if err != nil {
				t.Fatal(err)
			}
			got := packetEvaluate(p)
			if got != c.wantEV {
				t.Errorf("packetEvaluate()=%d, want %d", got, c.wantEV)
			}
		})
	}
}
