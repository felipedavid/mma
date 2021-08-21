package main

import (
    "strconv"
    "fmt"
)

func parseInteger(integer string) (value int, err error) {
    var tmp int64
    if tmp, err = strconv.ParseInt(integer, 10, 16); err == nil {
        value = int(tmp)
        return
    }

    if len(integer) > 2 {
        if tmp, err = strconv.ParseInt(integer[2:], 16, 16); err == nil {
            value = int(tmp)
            return
        }
    }
    fmt.Println("he");

    return
}

func isStringInt(str string) bool {
    if _, err := strconv.Atoi(str); err == nil {
        return true
    }
    return false
}

func isMemoryReference(str string) bool {
    return str[0] == '[' && str[len(str)-1] == ']'
}
