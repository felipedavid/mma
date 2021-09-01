package main

import "fmt"

type DataInterface interface {
    HexString() string
}

type DataFile struct {
    DataStream []DataInterface
}

type Data struct {
    lit string
    byte_data uint16
}

func (d *Data) HexString() string {
    return fmt.Sprintf("%04x", d.byte_data)
}
