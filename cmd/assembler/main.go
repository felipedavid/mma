package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func prompt() {
	asm := newAssembler("<string>", nil)
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
		asm.resetState([]byte(line))
		asm.parseLine()
		if !asm.hasError {
			fmt.Println(asm.getDebugStr())
		}
	}
}

func compile(fileName string) {
	source, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	asm := newAssembler(fileName, source)
	asm.parseLines()
}

func main() {
	interactiveMode := *flag.Bool("i", true, "Interactive mode")
	flag.Parse()

	if interactiveMode {
		prompt()
	} else {
		fileName := os.Args[1]
		compile(fileName)
	}
}
