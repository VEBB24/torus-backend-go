#!/bin/bash

go get "gopkg.in/resty.v0"
go get "github.com/gorilla/mux"

mkdir -p build

go build -o build/server src/*.go