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
    op    int
    rs    int
    rt    int
    rd    int
    funct int
}

func (r *RInstruction) BinaryString() string {
    binn := 0x0000
    binn |= r.op << 12
    binn |= r.rs << 9
    binn |= r.rt << 6
    binn |= r.rd << 3
    binn |= r.funct

    return fmt.Sprintf("%x", binn)
}

type IInstruction struct {
    lit string
    op   int
    rs   int
    rt   int
    immd int
}

func (i *IInstruction) BinaryString() string {
    binn := 0x0000;
    binn |= i.op << 12
    binn |= i.rs << 9
    binn |= i.rt << 6
    binn |= i.immd

    return fmt.Sprintf("%x", binn)
}
