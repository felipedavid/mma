package main

import "fmt"

func main() {
	asm := newAssembler("<string>", []byte("lw $2, 1($1)\n lw $3, 2($0)  \n add $4, $2, $3\n"))
	//for !asm.isAtEnd() {
	//	fmt.Printf("[kind: %s] [val: %v]\n", kindStr[asm.token.kind], asm.token.val)
	//	asm.nextToken()
	//}
	//fmt.Printf("[kind: %s] [val: %v]\n", kindStr[asm.token.kind], asm.token.val)
	asm.parseLine()
	asm.parseLine()
	fmt.Printf(asm.getCodeSectionStr())
}
