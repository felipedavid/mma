package main

import (
    "strconv"
    "os"
    "fmt"
)

func parseImmediate(immd string) (value int, err error) {
    var tmp int64
    if len(immd) > 2 {
    } else {
        tmp, err = strconv.ParseInt(immd, 10, 16)
    }
    value = int(tmp)

    if err != nil {
        fmt.Printf("Imediato '%v' invalido.\n", immd)
        os.Exit(1)
    }
    return
}
