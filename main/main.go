package main

import (
	"github.com/cnblvr/websocket/chat"
	"log"
	"net/http"
)

func main() {
	// server
	var server = chat.NewServer()
	go server.Listen()

	// client
	http.Handle("/", http.FileServer(http.Dir("webroot")))
	log.Panic(http.ListenAndServe(":8080", nil))
}
