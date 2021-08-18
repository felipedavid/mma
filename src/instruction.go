package main

import "fmt"

type DataInterface interface {
    BinaryString() string
}

type DataFile struct {
    DataStream []DataInterface
}

type Data struct {
    lit string
    byte_data uint16
}

func (d *Data) BinaryString() string {
    return fmt.Sprintf("%04x", d.byte_data)
}

type Instruction interface {
    BinaryString() string
}

type AssemblyFile struct {
    Instructions []Instruction
}

type RInstruction struct {
    lit string
    rs    uint16
    rt    uint16
    rd    uint16
    funct uint16
}

func (r *RInstruction) BinaryString() string {
    var bin uint16
    bin |= r.rs << 9
    bin |= r.rt << 6
    bin |= r.rd << 3
    bin |= r.funct
    return fmt.Sprintf("%04x", bin)
}

type IInstruction struct {
    lit string
    op   uint16
    rs   uint16
    rt   uint16
    immd uint16
}

func (i *IInstruction) BinaryString() string {
    var bin uint16
    bin |= i.op << 12
    bin |= i.rs << 9
    bin |= i.rt << 6
    bin |= i.immd
    return fmt.Sprintf("%04x", bin)
}

type JInstruction struct {
    lit string
    addr uint16
}

func (j *JInstruction) BinaryString() string {
    var bin uint16
    bin |= j.addr
    return fmt.Sprintf("%04x", bin)
}
