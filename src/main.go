package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
    "flag"
    "errors"
    "path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: ./assembler program.m")
        os.Exit(1)
	}

    flag.Parse()
    path := flag.Arg(0)
    asmData, err := readFile(path)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[!] Falha ao tentar ler arquivo: %s\n", err)
        os.Exit(1)
    }

	var p Parser
	p.Init(asmData)
	dataFile, assemblyFile := p.Parse()

    if len(dataFile.DataStream) > 0 {
        var dataBuffer bytes.Buffer
        dataBuffer.WriteString("v2.0 raw\n")
        for _, i := range dataFile.DataStream {
            dataBuffer.WriteString(i.HexString() + "\n")
        }
	    datFilePath := strings.Replace(path, ".m", ".dat", 1)
	    ioutil.WriteFile(datFilePath, dataBuffer.Bytes(), 0644)
    }

    if len(assemblyFile.Instructions) > 0 {
        var instructionBuffer bytes.Buffer
        instructionBuffer.WriteString("v2.0 raw\n")
        for _, i := range assemblyFile.Instructions {
            instructionBuffer.WriteString(i.HexString() + "\n")
            i.printDecode()
        }
        binFilePath := strings.Replace(path, ".m", ".ins", 1)
        ioutil.WriteFile(binFilePath, instructionBuffer.Bytes(), 0644)
    }
}

func readFile(path string) ([]byte, error) {
    if path == "" {
        return []byte{}, errors.New("selecione um arquivo \".m\"")
    }

    if filepath.Ext(path) != ".m" {
        return []byte{}, errors.New("selecione um arquivo \".m\"")
    }

    bytes, err := os.ReadFile(path)
    if err != nil {
        return []byte{}, errors.New("erro ao ler arquivo")
    }
    return bytes, nil
}
