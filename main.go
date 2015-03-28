package main

import (
	"net/http"
	"os"

	"github.com/op/go-logging"
	"golang.org/x/net/websocket"
)

var log = logging.MustGetLogger("blocky")

const (
	logFormat = "%{time:15:04:05.000} %{level:.4s} [%{shortfunc}] %{message}"
	listen    = ":1987"
)

func client(ws *websocket.Conn) {
}

func region(id string) {
}

func main() {
	// Set up logging.
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, logging.MustStringFormatter(logFormat))
	logging.SetBackend(formatter)

	// Start the websocket server.
	log.Info("Listening on %s...", listen)
	http.Handle("/socket", websocket.Handler(client))
	log.Fatal(http.ListenAndServe(listen, nil))
}
