package main

import (
	"testing"
)

func TestAssembler(t *testing.T) {
	tests := []struct {
		name  string
		instr string
		want  uint16
	}{
		{
			name:  "lw",
			instr: "lw $2, 0($1)",
			want:  0b0011_001_010_000000,
		},
		{
			name:  "add",
			instr: "add $4, $2, $3",
			want:  0b0000_010_011_100_000,
		},
		{
			name:  "sw",
			instr: "sw $4, 4($1)",
			want:  0b1010_001_100_000100,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			asm := newAssembler("<string>", []byte(test.instr))
			asm.parseLine()

			high := uint16(asm.codeSection[0])
			low := uint16(asm.codeSection[1])

			mc := low | (high << 8)

			if mc != test.want {
				t.Errorf("got %016b; want %016b", mc, test.want)
			}
		})
	}
}
