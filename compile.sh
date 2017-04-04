#!/bin/bash

go get "gopkg.in/resty.v0"
go get "github.com/gorilla/mux"
go get "github.com/colinmarc/hdfs"
go get "github.com/mediocregopher/radix.v2/redis"
go get "github.com/golang/glog"

mkdir -p build

go build -o build/server src/*.go

docker-compose build