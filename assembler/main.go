package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: ./%s <program>\n", os.Args[0])
		os.Exit(-1)
	}

	source, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	a := newAssembler(source)
	err = a.assemble()
	if err != nil {
		panic(err)
	}
}
