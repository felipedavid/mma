package main

import "fmt"

func main() {
	lex := NewLexer("no_name", []byte(" 123123 \"hello\" 23 76\n 0b502 0xff"))
	for lex.Token.kind != End {
		fmt.Printf("[TokenType: %s] [TokenVal: %v]\n", TokenKindToString[lex.Token.kind], lex.Token.GetValue())
		lex.NextToken()
	}
}
