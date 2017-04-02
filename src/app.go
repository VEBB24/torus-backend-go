package main

import (
	"log"
	"net/http"
)

func main() {
	r := NewRouter()
	log.Fatal(http.ListenAndServe(":12345", r))
}
