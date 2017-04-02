package main

import (
	"log"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	r := NewRouter()

	glog.Infoln("Start server")
	log.Fatal(http.ListenAndServe(":12345", r))
}
