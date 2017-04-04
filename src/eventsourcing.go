package main

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

//User is keeping token and the message channel
type User struct {
	// User token
	token string

	// Channel to push message
	channel chan string
}

//Message is keeping the user and a custom message
type Message struct {
	// User token destination
	to string

	// String representing a message
	msg string
}

//Broker is keeping a list of current client
type Broker struct {
	// Create a map of clients
	clients map[string](chan string)

	// Channel into which new clients can be pushed
	newClients chan *User

	// Channel into which disconnected clients should be pushed
	defunctClients chan *User

	// Channel into which messages are pushed to be broadcast out
	// to attahed clients.
	messages chan *Message
}

//Start new goroutine for handling the add & removal fo clients,
//and broadcast messages
func (b *Broker) Start() {

	go func() {

		for {

			select {

			case user := <-b.newClients:

				b.clients[user.token] = user.channel
				glog.Infoln("Added new client")

			case user := <-b.defunctClients:

				delete(b.clients, user.token)
				close(user.channel)

				glog.Infoln("Removed client")

			case message := <-b.messages:
				receiver := message.to
				if _, ok := b.clients[receiver]; ok {
					b.clients[receiver] <- message.msg
					glog.Infof("Send message to %s clients", receiver)
				} else {
					glog.Infof("Client does not exist")
				}
			}
		}
	}()
}

// ServeHTTP method handles and HTTP request
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Make sure that the writer supports flushing.
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)

	user := &User{
		token:   params["id"],
		channel: make(chan string),
	}

	// Add this client to the map of those that should
	// receive updates
	b.newClients <- user

	// Listen to the closing of the http connection via the CloseNotifier
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		//remove client
		b.defunctClients <- user
		glog.Infoln("HTTP connection just closed.")
	}()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {

		// Read from user channel.
		msg, open := <-user.channel

		if !open {
			// If our user.channel was closed, this means that the client has
			// disconnected.
			break
		}

		// Write to the ResponseWriter, `w`.
		fmt.Fprintf(w, "data: Message: %s\n\n", msg)

		f.Flush()
	}

	glog.Infoln("Finished HTTP request at ", r.URL.Path)
}
