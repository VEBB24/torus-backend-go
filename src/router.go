package main

import "github.com/gorilla/mux"
import "github.com/golang/glog"

//NewRouter create a new router instance
func NewRouter() *mux.Router {
	glog.Infoln("Declare new router")

	router := mux.NewRouter()
	router.HandleFunc("/auth", checkAuth).Methods("POST")
	return router
}
