#!/bin/bash

go get "gopkg.in/resty.v0"
go get "github.com/gorilla/mux"
go get "github.com/colinmarc/hdfs"

mkdir -p build

go build -o build/server src/*.go