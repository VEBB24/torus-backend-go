package main

import (
	"log"
	"net/http"

	"flag"

	"path/filepath"

	"github.com/golang/glog"
)

var basePath *string
var baseHost *string
var redisClient *Redis

func main() {
	baseHost = flag.String("host", "127.0.0.1", "a valid host")
	port := flag.String("port", "12345", "a valid port")
	path := flag.String("basePath", ".", "a valid path")

	flag.Parse()

	tmp, _ := filepath.Abs(*path)
	basePath = &tmp

	redisClient = RedisFactory(*baseHost, 10)

	glog.Infoln("Start server")
	log.Fatal(http.ListenAndServe(*baseHost+":"+*port, NewRouter()))
}
