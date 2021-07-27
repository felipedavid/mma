#!/usr/bin/python3
# Observação: Quase nunca uso python e queria acabar isso o mais rápido possível..
# Então pode esperar um código de ótima "qualidade", haha :v

import sys
import re

type_r = {"ADD":0, "AND":4, "NAND":1, "NOR":7, "OR":5, "SLT":2, "SUB":3}
type_i = {"ADDI":8, "BEQ":4, "LW":3, "SW":0xa}

def main():
    if len(sys.argv) != 2:
        print("Usage: ./assembler program.m")
        return

    f = open(sys.argv[1], "r")
    
    addr = 0
    for i in f:
        tok = list(filter(None, re.split(",| ", i)))
        iname = tok[0].upper()
        machine_i = ""

        print("{0:#0{1}x}".format(addr,6), " ", end='')
        #print(i.strip(), "-> ", end="")
        print(tok, "-> ", end="")
        if iname in type_r.keys():
            print(encode_r_instruction(tok, iname))
        elif iname in type_i:
            print(encode_i_instruction(tok, iname))
        elif iname == "J":
            print(encode_j_instruction(tok))
        else:
            printf("ERROR:" + iname + "is not a valid instruction.\n")
            return
        addr += 2

def bin_string(value: int, padding: int) -> str:
    return format(int(value), '0'+str(padding)+'b')

def encode_r_instruction(tokens: list, iname: str) -> str:
    opc = bin_string(0, 4)
    rs = bin_string(tokens[2][1], 3)
    rt = bin_string(tokens[3][1], 3)
    rd = bin_string(tokens[1][1], 3)
    funct = bin_string(type_r[iname], 3)

    bin_instruction = opc + rs + rt + rd + funct
    return '0x' + format(int(bin_instruction, 2), '04x')
    
def encode_i_instruction(tokens: list, iname: str) -> str:
    return "0x0000"

def encode_j_instruction(tokens: list) -> str:
    opc = bin_string(2, 4)
    immediate = bin_string(int(tokens[1].strip(), 0), 12)
    bin_instruction = opc + immediate
    return '0x' + format(int(bin_instruction, 2), '04x')

if __name__ == "__main__":
    main()
