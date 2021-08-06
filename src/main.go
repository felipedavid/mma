package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "bytes"
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
    assemblyFile := p.Parse()

    var b bytes.Buffer
    for _, i := range assemblyFile.Instructions {
        b.WriteString(i.BinaryString())
    }

    executablePath := strings.Replace(asmFilePath, ".m", ".i", 1)
    ioutil.WriteFile(executablePath, b.Bytes(), 0644)
}
