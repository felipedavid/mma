package main

import "fmt"

func main() {
	asm := newAssembler("<string>", []byte("add $1, $2, $3,"))
	for !asm.isAtEnd() {
		fmt.Printf("[kind: %s] [val: %v]\n", kindStr[asm.token.kind], asm.token.val)
		asm.nextToken()
	}
	fmt.Printf("[kind: %s] [val: %v]\n", kindStr[asm.token.kind], asm.token.val)
}
