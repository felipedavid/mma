package main

import "fmt"

func main() {
	lex := NewLexer([]byte(" 123123  23 76 0b101 0xff"))
	for lex.Token.kind != End {
		fmt.Printf("[TokenType: %d] [TokenVal: %d]\n", lex.Token.kind, lex.Token.intVal)
		lex.NextToken()
	}
}
