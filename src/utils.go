package main

import (
    "strconv"
)

func parseInteger(integer string) (value int, err error) {
    var tmp int64
    if tmp, err = strconv.ParseInt(integer, 10, 16); err == nil {
        value = int(tmp)
        return
    }

    if len(integer) > 3 {
        if tmp, err = strconv.ParseInt(integer[2:], 16, 16); err == nil {
            value = int(tmp)
            return
        }
    }

    return
}

func isStringInt(str string) bool {
    if _, err := strconv.Atoi(str); err == nil {
        return true
    }
    return false
}
