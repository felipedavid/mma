package main

import "fmt"

func main() {
	asm := newAssembler("<string>", []byte("abcefg \"hello there\"  "))
	for !asm.isAtEnd() {
		fmt.Printf("[kind: %d] [val: %v]\n", asm.token.kind, asm.token.val)
		asm.nextToken()
	}
}
