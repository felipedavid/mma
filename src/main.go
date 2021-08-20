package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./assembler program.m")
		os.Exit(1)
	}
	asmFilePath := os.Args[1]
	asmData, _ := ioutil.ReadFile(asmFilePath)

	var p Parser
	p.Init(asmData)
	dataFile, assemblyFile := p.Parse()

    if len(dataFile.DataStream) > 0{
        var d bytes.Buffer
        d.WriteString("v2.0 raw\n")
        for _, i := range dataFile.DataStream {
            d.WriteString(i.BinaryString() + "\n")
        }
	    datFilePath := strings.Replace(asmFilePath, ".m", ".dat", 1)
	    ioutil.WriteFile(datFilePath, d.Bytes(), 0644)
    }

    if len(assemblyFile.Instructions) > 0{
        var b bytes.Buffer
        b.WriteString("v2.0 raw\n")
        for _, i := range assemblyFile.Instructions {
            b.WriteString(i.BinaryString() + "\n")
        }
        binFilePath := strings.Replace(asmFilePath, ".m", ".ins", 1)
        ioutil.WriteFile(binFilePath, b.Bytes(), 0644)
    }
}
