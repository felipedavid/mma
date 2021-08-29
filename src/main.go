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

    if len(dataFile.DataStream) > 0 {
        var dataBuffer bytes.Buffer
        dataBuffer.WriteString("v2.0 raw\n")
        for _, i := range dataFile.DataStream {
            dataBuffer.WriteString(i.HexString() + "\n")
        }
	    datFilePath := strings.Replace(asmFilePath, ".m", ".dat", 1)
	    ioutil.WriteFile(datFilePath, dataBuffer.Bytes(), 0644)
    }

    if len(assemblyFile.Instructions) > 0 {
        var instructionBuffer bytes.Buffer
        instructionBuffer.WriteString("v2.0 raw\n")
        for _, i := range assemblyFile.Instructions {
            instructionBuffer.WriteString(i.HexString() + "\n")
            i.printDecode()
        }
        binFilePath := strings.Replace(asmFilePath, ".m", ".ins", 1)
        ioutil.WriteFile(binFilePath, instructionBuffer.Bytes(), 0644)
    }
}
