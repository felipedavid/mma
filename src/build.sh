#!/bin/bash

bin_path=../bin/

if [[ ! -d $bin_path ]]; then
    mkdir $bin_path
fi

go build -o $bin_path .
GOOS=windows GOARCH=amd64 go build -o $bin_path .
