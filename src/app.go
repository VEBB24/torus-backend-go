package main

import (
	"log"
	"net/http"

	"flag"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	r := NewRouter()

	glog.Infoln("Start server")
	log.Fatal(http.ListenAndServe(":12345", r))
}
