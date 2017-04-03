package main

import (
	"log"
	"net/http"

	"flag"

	"path/filepath"

	"github.com/golang/glog"
)

var basePath *string
var redisClient *Redis

func main() {
	host := flag.String("host", "0.0.0.0", "a valid host")
	port := flag.String("port", "12345", "a valid port")
	path := flag.String("basePath", ".", "a valid path")
	flag.Parse()
	tmp, _ := filepath.Abs(*path)
	basePath = &tmp
	redisClient = RedisFactory(*host)
	r := NewRouter()
	glog.Infoln("Start server")
	log.Fatal(http.ListenAndServe(*host+":"+*port, r))
}
