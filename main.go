package main

import (
	"net/http"
	"os"

	"github.com/blixt/geomys"
	"github.com/op/go-logging"

	"./blocky"
)

var (
	log      = logging.MustGetLogger("blocky")
	universe = blocky.NewUniverse()
)

const (
	logFormat = "%{time:15:04:05.000} %{level:.4s} [%{longfunc}] %{message}"
	listen    = ":1987"
)

func main() {
	// Set up logging.
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, logging.MustStringFormatter(logFormat))
	logging.SetBackend(formatter)

	// Run the universe (no big deal).
	go universe.Run()

	// Start the websocket server.
	log.Info("Starting server on %s...", listen)
	http.Handle("/socket", geomys.WebSocketHandler(universe))
	log.Fatal(http.ListenAndServe(listen, nil))
}
