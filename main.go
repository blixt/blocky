package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func client(ws *websocket.Conn) {
}

func region(id string) {
}

func main() {

	// Start the websocket server.
	http.Handle("/socket", websocket.Handler(client))
	log.Fatal(http.ListenAndServe(":12345", nil))
}
