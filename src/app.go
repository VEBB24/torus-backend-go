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
var baseRedis *string
var baseHdfs *string

func main() {
	baseHost = flag.String("host", "127.0.0.1", "a valid host")
	baseRedis = flag.String("redis", "127.0.0.1", "a valid redis host")
	baseHdfs = flag.String("hdfs", "127.0.0.1", "a valid hdfs host")
	port := flag.String("port", "12345", "a valid port")
	path := flag.String("basePath", ".", "a valid path")

	flag.Parse()

	tmp, _ := filepath.Abs(*path)
	basePath = &tmp

	redisClient = RedisFactory(*baseRedis, 10)

	glog.Infoln("Start server")
	log.Fatal(http.ListenAndServe(":"+*port, NewRouter()))
}
