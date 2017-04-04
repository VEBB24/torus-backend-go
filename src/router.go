package main

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

var broker *Broker

func init() {
	broker = &Broker{
		make(map[string](chan string)),
		make(chan *User),
		make(chan *User),
		make(chan *Message),
	}
	broker.Start()
}

//NewRouter create a new router instance
func NewRouter() *mux.Router {
	glog.Infoln("Declare new router")

	router := mux.NewRouter()
	router.HandleFunc("/auth", checkAuth).Methods("POST")
	router.HandleFunc("/hdfs/{id}", getFiles).Methods("GET")
	router.HandleFunc("/hdfs/{id}/{file}", removeFile).Methods("DELETE")
	router.HandleFunc("/hdfs/{id}", renameFile).Methods("PUT")
	router.Handle("/streaming/{id}", broker).Methods("GET")

	go func() {
		for i := 0; ; i++ {
			m := &Message{
				msg: fmt.Sprintf("%d - the time is %v", i, time.Now()),
				to:  "1",
			}
			broker.messages <- m
			time.Sleep(3 * 1e9)
		}
	}()

	return router
}
