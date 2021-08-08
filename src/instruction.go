package main

import (
    "fmt"
)

type Instruction interface {
    BinaryString() string
}

type AssemblyFile struct {
    Instructions []Instruction
}

type RInstruction struct {
    lit string
    rs    int
    rt    int
    rd    int
    funct int
}

func (r *RInstruction) BinaryString() string {
    bin := 0x0000
    bin |= r.rs << 9
    bin |= r.rt << 6
    bin |= r.rd << 3
    bin |= r.funct
    return fmt.Sprintf("%04x", bin)
}

type IInstruction struct {
    lit string
    op   int
    rs   int
    rt   int
    immd int
}

func (i *IInstruction) BinaryString() string {
    bin := 0x0000
    bin |= i.op << 12
    bin |= i.rs << 9
    bin |= i.rt << 6
    bin |= i.immd
    return fmt.Sprintf("%04x", bin)
}

type JInstruction struct {
    lit string
    addr int
}

func (j *JInstruction) BinaryString() string {
    bin := 0x2000
    bin |= j.addr
    return fmt.Sprintf("%04x", bin)
}
