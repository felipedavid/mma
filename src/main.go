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
    dataFile, assemblyFile := p.Parse()

    var d bytes.Buffer
    for n, i := range dataFile.DataStream {
        if n == 0 {
            d.WriteString("v2.0 raw\n")
        }
        d.WriteString(i.BinaryString() + "\n")
    }

    var b bytes.Buffer
    for n, i := range assemblyFile.Instructions {
        if n == 0 {
            b.WriteString("v2.0 raw\n")
        }
        b.WriteString(i.BinaryString() + "\n")
    }

    datFilePath := strings.Replace(asmFilePath, ".m", ".dat", 1)
    binFilePath := strings.Replace(asmFilePath, ".m", ".ins", 1)

    ioutil.WriteFile(datFilePath, d.Bytes(), 0644)
    ioutil.WriteFile(binFilePath, b.Bytes(), 0644)
}
