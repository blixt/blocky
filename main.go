package main

import (
	"net/http"
	"os"

	"github.com/op/go-logging"
	"golang.org/x/net/websocket"

	"./blocky"
)

var (
	context = blocky.Context{Version: "0.1.0.001"}
	log     = logging.MustGetLogger("blocky")
)

const (
	logFormat = "%{time:15:04:05.000} %{level:.4s} [%{shortfunc}] %{message}"
	listen    = ":1987"
)

func client(ws *websocket.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Error in client: %s", r)
		}
	}()

	// Shake hands.
	var hello *blocky.Hello
	websocket.JSON.Receive(ws, &hello)
	websocket.JSON.Send(ws, blocky.Handshake(context, hello))

	//
}

func region(id string) {
}

func main() {
	// Set up logging.
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, logging.MustStringFormatter(logFormat))
	logging.SetBackend(formatter)

	// Start the websocket server.
	log.Info("Starting server on %s...", listen)
	http.Handle("/socket", websocket.Handler(client))
	log.Fatal(http.ListenAndServe(listen, nil))
}
