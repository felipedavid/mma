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
    debug := flag.Bool("d", false, "executar em modo debugging")

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

    // Se flag de debugging está ativada, printar converções e sair
    if *debug && len(assemblyFile.Instructions) > 0 {
        for _, i := range assemblyFile.Instructions {
            i.printDecode()
        }
        fmt.Println()
        return
    }

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
        }
        binFilePath := strings.Replace(path, ".m", ".ins", 1)
        ioutil.WriteFile(binFilePath, instructionBuffer.Bytes(), 0644)
    }
}

func readFile(path string) ([]byte, error) {
    if path == "" || filepath.Ext(path) != ".m" {
        return []byte{}, errors.New("selecione um arquivo \".m\"")
    }

    bytes, err := os.ReadFile(path)
    if err != nil {
        return []byte{}, errors.New("erro ao ler arquivo")
    }
    return bytes, nil
}
