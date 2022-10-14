package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	//interactiveMode := *flag.Bool("i", false, "Interactive mode")
	flag.Parse()

	asm := newAssembler("<string>", []byte("sw $4, 4($1)"))
	asm.parseLines()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">>> ")
		line, _ := reader.ReadString('\n')
		if line == "" {
			break
		} else if line == "\n" {
			continue
		}
		line = strings.Replace(line, "\n", "", -1)
		asm.setSource([]byte(line))
		asm.parseLines()
		fmt.Println(asm.getDebugStr())
	}
}
