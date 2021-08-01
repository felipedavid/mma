#!/usr/bin/python3

import sys
import re

type_r = {
    "ADD":  0x0000, 
    "AND":  0x0004, 
    "NAND": 0x0001, 
    "NOR":  0x0007, 
    "OR":   0x0005,  
    "SLT":  0x0002, 
    "SUB":  0x0003,
}

type_i = {
    "ADDI": 0x8000, 
    "BEQ":  0x4000, 
    "LW":   0x3000, 
    "SW":   0xa000,
}

labels = {}

def main():
    if len(sys.argv) != 2:
        print("Usage: ./assembler program.m")
        return

    as_file = open(sys.argv[1], "r")
    bi_file = open("a.out", "w")

    ibuffer = []

    addr = 0
    for i in as_file:
        tok = list(filter(None, re.split(",| ", i)))
        tok[-1] = tok[-1].strip()
        iname = tok[0].upper()

        if len(tok) == 1 and tok[0][-1] == ':':
            tok[0] = tok[0].replace(":", "")
            labels[tok[0]] = addr + 2
        elif iname in type_r:
            ibuffer.append(encode_r_instruction(tok, iname, addr))
        elif iname in type_i:
            ibuffer.append(encode_i_instruction(tok, iname, addr))
        elif iname == "J":
            ibuffer.append(encode_j_instruction(tok, iname, addr))

        addr += 2

    bi_file.write("v2.0 raw\n")
    for i in ibuffer:
        bi_file.write(format(i, '04x') + '\n')
    

def encode_r_instruction(tokens: list, iname: str, addr: int) -> int:
    rs = int(tokens[2][1])
    rt = int(tokens[3][1])
    rd = int(tokens[1][1])

    mi = type_r[iname]
    v_rang = range(0, 7)

    if rs in v_rang and rt in v_rang and rd in v_rang:
        mi |= (rs << 9)
        mi |= (rt << 6)
        mi |= (rd << 3)
    else:
        print(f"ERRO: Linha {addr}. Registrador inválido")
        exit()

    return mi;

def encode_i_instruction(tokens: list, iname: str, addr: int) -> int:
    if iname == "BEQ":
        rt = int(tokens[2][1])
        rs = int(tokens[1][1])
        immd = int(tokens[3].strip(), 0)
    else:
        rt = int(tokens[1][1])
        rs = int(tokens[2].split("$", 1)[1][0])
        immd = int(tokens[2].split("$", 1)[0].replace("(", ""))

    mi = type_i[iname]
    mi |= (rs << 9)
    mi |= (rt << 6)
    mi |= immd

    return mi


def encode_j_instruction(tokens: list, iname: str, addr: int) -> int:
    if (tokens[1][0].isalpha()):
        tokens[1] = labels[tokens[1]]
    else:
        tokens[1] = tokens[1].strip()
    return 0x2000 | tokens[1]

if __name__ == "__main__":
    main()
