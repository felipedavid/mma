package main

import (
    "strings"
    "strconv"
    "fmt"
    "os"
)

type Parser struct {
    scanner Scanner
    symbols map[string]int
    instructions []Instruction
    dataStream []DataInterface
    nAddr int
}

func (p *Parser) Init(src []byte) {
    p.scanner.Init(src)
    p.nAddr = 0
    p.symbols = map[string]int {
        "R0": 0, "R1": 1, "R2": 2, "R3": 3, "R4": 4, "R5": 5, "R6": 6, "R7": 7,
    }
}

const (
    INSTRUCTION = 0
    DATA = 1
)

func (p *Parser) appendPseudoInstruction(lit string) {
    no_commas := strings.ReplaceAll(lit, ",", " ")
    lits := strings.Fields(no_commas)

    var instruction string
    switch lits[0] {
    case "li":
        to := lits[1]
        from := lits[2]
        if p.isRegister(to) && isStringInt(from) {
            instruction = "addi " + to + ", $r0, " + from
        } else {
            fmt.Printf("Formato da instrução \"%v\" é inválido.\n", lit)
            os.Exit(1)
        }
        p.instructions = append(p.instructions, &IInstruction{lit: instruction})
    }
}

func (p *Parser) Parse() (DataFile, AssemblyFile) {
    ins_or_data := INSTRUCTION
loop:
    for {
        tok, lit := p.scanner.Scan()
        if (lit == ".data") {
            ins_or_data = DATA
            p.symbols[lit] = len(p.instructions)
            continue
        } else if lit == ".text" {
            ins_or_data = INSTRUCTION
            p.symbols[lit] = len(p.instructions)
            continue
        }
        lit = strings.TrimSpace(lit)

        if ins_or_data == INSTRUCTION {
            switch tok {
            case EOF:
                break loop
            case LABEL:
                if _, ok := p.symbols[lit]; ok {
                    fmt.Printf(`[!] Label \"%v\" já foi declarada. \n[!] 
                        Não é possível definir múltiplas labels com o mesmo nome.\n`, lit)
                    os.Exit(1)
                }
                p.symbols[lit[:len(lit)-1]] = len(p.instructions)
            case R_INSTRUCTION:
                p.instructions = append(p.instructions, &RInstruction{lit: lit})
            case I_INSTRUCTION:
                p.instructions = append(p.instructions, &IInstruction{lit: lit})
            case J_INSTRUCTION:
                p.instructions = append(p.instructions, &JInstruction{lit: lit})
            case PSEUDO_INSTRUCTION:
                p.appendPseudoInstruction(lit)
            }
        } else if ins_or_data == DATA {

            switch tok{
            case LABEL:
                lits := strings.Fields(lit)
                label_name := lits[0][:len(lits[0])-1]
                if _, ok := p.symbols[label_name]; ok {
                    fmt.Printf(`[!] Label \"%v\" já foi declarada.\n
                    [!] Não é possível definir múltiplas labels com o mesmo nome.\n`, label_name)
                    os.Exit(1)
                }
                p.symbols[label_name] = len(p.dataStream) + 1
                if len(lits) >= 3 {
                    lits = lits[1:]
                    p.dataStream = append(p.dataStream, &Data{lit: lits[0] + " " + lits[1]})
                }
            case WORD:
                p.dataStream = append(p.dataStream, &Data{lit: lit})
            }
        }
    }

    for _, dat := range p.dataStream {
        switch d := dat.(type) {
        case *Data:
            p.parseData(d)
        }
    }

    for _, instr := range p.instructions {
        switch i := instr.(type) {
        case *RInstruction:
            p.parseRInstruction(i)
        case *IInstruction:
            p.parseIInstruction(i)
        case *JInstruction:
            p.parseJInstruction(i)
        }
    }
    return DataFile{DataStream: p.dataStream}, AssemblyFile{ Instructions: p.instructions }
}

func (p *Parser) parseData(d *Data) {
    lits := strings.Fields(d.lit)

    if lits[0] == ".word" {
        n, err := parseInteger(lits[1])
        if err != nil {
            fmt.Println("Declaração de dado inválida: ", d.lit)
            os.Exit(1)
        }
        d.byte_data = uint16(n)
        return
    } else {
        fmt.Printf("Tipo \"%v\" não suportado. No momento apenas o tipo \".word\" é suportado", lits[0])
        os.Exit(1)
    }
}

func (p *Parser) parseRInstruction(r *RInstruction) {
    no_commas := strings.ReplaceAll(r.lit, ",", " ")
    lits := strings.Fields(no_commas)

    // Parsing funct
    switch lits[0] {
        case "add":
            r.funct = 0
        case "and":
            r.funct = 4
        case "nand":
            r.funct = 1
        case "nor":
            r.funct = 7
        case "or":
            r.funct = 5
        case "slt":
            r.funct = 2
        case "sub":
            r.funct = 3
        default:
            fmt.Printf("[!] ERRO: Instrução inválida: %v", r.lit)
            os.Exit(1)
    }

    var registers []uint16
    registers, lits = p.parseRegisters(lits)
    if len(registers) != 3 {
        fmt.Printf("[!] ERRO: Instrução \"%v\" deveria referenciar 3 registradores.\n", r.lit)
        os.Exit(1)
    }
    r.rd = registers[0]
    r.rs = registers[1]
    r.rt = registers[2]
}

func (p *Parser) parseIInstruction(i *IInstruction) {
    clean_lit := strings.ReplaceAll(i.lit, ",", " ")
    clean_lit = strings.ReplaceAll(clean_lit, "(", " ")
    clean_lit = strings.ReplaceAll(clean_lit, ")", " ")
    lits := strings.Fields(clean_lit)

    switch lits[0] {
        case "addi":
            i.op = 8
        case "beq":
            i.op = 4
        case "lw":
            i.op = 3
        case "sw":
            i.op = 10
        default:
            fmt.Println("Invalid instruction:", i.lit)
            os.Exit(1)
    }

    var registers []uint16
    registers, lits = p.parseRegisters(lits)
    if len(registers) == 2 && len(lits) == 1 && isStringInt(lits[0]) {
        if i.op == 3 || i.op == 10 || i.op == 8 {
            i.rt = registers[0]
            i.rs = registers[1]
        } else if (i.op == 4) {
            i.rt = registers[1]
            i.rs = registers[0]
        } else {
            fmt.Printf("[!] Instrução \"%v\" inválida.", i.lit)
            os.Exit(1)
        }
        val, _ := parseInteger(lits[0])
        i.immd = uint16(val)
    } else if len(registers) == 1 && len(lits) == 1 && (i.op == 3 || i.op == 10) {
        if val, ok := p.symbols[lits[0]]; ok {
            i.immd = uint16(val)
            i.rt = registers[0]
            i.rs = 0
        } else {
            fmt.Printf("[!] Label \"%v\" não existe.\n", lits[0])
            os.Exit(1)
        }
    } else {
        fmt.Printf("[!] Instrução \"%v\" inválida.\n", i.lit)
        os.Exit(1)
    }
}

func (p *Parser) parseJInstruction(j *JInstruction) {
    label := strings.Fields(j.lit)[1]
    if val, ok := p.symbols[label]; ok { // if the label is present in the map
        j.addr = uint16(val)
    } else {
        addr, err := parseInteger(label)
        if err != nil {
            fmt.Printf("Label \"%v\" não existe.\n", label)
            os.Exit(1)
        }
        j.addr = uint16(addr) / 2
    }
}

func (p *Parser) isLabel(lit string) bool {
    if _, ok := p.symbols[lit]; ok {
        return true
    }

    if len(lit) > 2 {
        if _, ok := p.symbols[lit[1:len(lit)-1]]; ok {
            return true
        }
    }

    return false
}

func (p *Parser) isRegister(str string) bool {
    if str[0] != '$' {
        return false
    }
    str = str[1:]

    if isStringInt(str) {
        return true
    }

    if _, ok := p.symbols[strings.ToUpper(str)]; ok {
        return true
    }

    fmt.Printf("Registrador \"%v\" não existe.\n", str)
    return false
}

func (p *Parser) parseRegisters(lit []string) (registers []uint16, new_lit []string) {
    var register_id uint16
    var tmp int

    lit = lit[1:]
    for _, v := range lit {
        register_name := v[1:]
        if p.isRegister(v) {
            if isStringInt(register_name) {
                tmp, _ = strconv.Atoi(register_name)
                register_id = uint16(tmp)
            } else {
                register_id = uint16(p.symbols[strings.ToUpper(v[1:])])
            }
            registers = append(registers, register_id)
        } else {
            new_lit = append(new_lit, v)
        }
    }
    return
}
