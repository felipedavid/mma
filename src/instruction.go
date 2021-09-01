package main

import "fmt"

type Instruction interface {
    HexString() string
    printDecode()
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

func (r *RInstruction) HexString() string {
    var bin uint16
    bin |= r.rs << 9
    bin |= r.rt << 6
    bin |= r.rd << 3
    bin |= r.funct
    return fmt.Sprintf("%04x", bin)
}

func (r *RInstruction) printDecode() {
    fmt.Printf("%-18v -> %04b %03b %03b %03b %03b -> 0x%v\n", r.lit, 0, r.rs, r.rt, r.rd, r.funct, r.HexString())
}

type IInstruction struct {
    lit string
    op   uint16
    rs   uint16
    rt   uint16
    immd uint16
}

func (i *IInstruction) HexString() string {
    var bin uint16
    bin |= i.op << 12
    bin |= i.rs << 9
    bin |= i.rt << 6
    bin |= i.immd
    return fmt.Sprintf("%04x", bin)
}

func (i *IInstruction) printDecode() {
    fmt.Printf("%-18v -> %04b %03b %03b %06b -> 0x%v\n", i.lit, i.op, i.rs, i.rt, i.immd, i.HexString())
}

type JInstruction struct {
    lit string
    addr uint16
}

func (j *JInstruction) HexString() string {
    var bin uint16
    bin |= j.addr
    bin |= 2 << 12;
    return fmt.Sprintf("%04x", bin)
}

func (j *JInstruction) printDecode() {
    fmt.Printf("%-18v -> %04b %012b -> 0x%v\n", j.lit, 2, j.addr, j.HexString())
}
