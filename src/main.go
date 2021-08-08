package main

import (
    "io/ioutil"
    "bytes"
    "os"
    "strings"
    "fmt"
)

func main() {
    if (len(os.Args) != 2) {
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
        fmt.Println(i.BinaryString())
        b.WriteString(i.BinaryString())
    }

    binFilePath := strings.Replace(asmFilePath, ".m", ".i", 1)
    ioutil.WriteFile(binFilePath, b.Bytes(), 0644)
}
