package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/op/go-logging"
	"golang.org/x/net/websocket"

	"./blocky"
)

var (
	context  = blocky.Context{Version: "0.1.0.001"}
	log      = logging.MustGetLogger("blocky")
	universe = blocky.NewUniverse()
)

const (
	logFormat = "%{time:15:04:05.000} %{level:.4s} [%{shortfunc}] %{message}"
	listen    = ":1987"
)

type intermediate struct {
	Type  string
	Value json.RawMessage
}

func receive(ws *websocket.Conn) (interface{}, error) {
	input := new(intermediate)
	if err := websocket.JSON.Receive(ws, input); err != nil {
		return nil, fmt.Errorf("Failed to receive packet from client (%s)", err)
	}

	var packet interface{}
	switch input.Type {
	case "Hello":
		packet = new(blocky.Hello)
	default:
		return nil, fmt.Errorf("Unsupported packet type %s", input.Type)
	}

	if err := json.Unmarshal(input.Value, packet); err != nil {
		return nil, fmt.Errorf("Failed to parse %s packet from client (%s)", input.Type, err)
	}

	return packet, nil
}

func mustReceive(ws *websocket.Conn) interface{} {
	if packet, err := receive(ws); err != nil {
		panic(err)
	} else {
		return packet
	}
}

func send(ws *websocket.Conn, value interface{}) error {
	valueType := reflect.TypeOf(value).Elem().Name()
	valueRaw, _ := json.Marshal(value)
	if err := websocket.JSON.Send(ws, &intermediate{valueType, valueRaw}); err != nil {
		return fmt.Errorf("Failed to send packet to client (%s)", err)
	}
	return nil
}

func mustSend(ws *websocket.Conn, value interface{}) {
	if err := send(ws, value); err != nil {
		panic(err)
	}
}

func disconnectClient(ws *websocket.Conn) {
	if r := recover(); r != nil {
		log.Error("Client error: %s", r)
	}
	send(ws, &blocky.Bye{"Closing connection"})
	if err := ws.Close(); err != nil {
		log.Error("Failed to close web socket (%s)", err)
	}
	log.Info("Client disconnected")
}

func receiveToInteractor(ws *websocket.Conn, ui *blocky.UniverseInteractor) {
	for {
		if packet, err := receive(ws); err != nil {
			log.Debug("Stopping receive: %s", err)
			ui.Close()
			break
		} else {
			log.Debug("Received packet %T", packet)
			ui.Put(packet)
		}
	}
}

func client(ws *websocket.Conn) {
	defer disconnectClient(ws)
	log.Info("Client connected")

	// Shake hands.
	hello := mustReceive(ws).(*blocky.Hello)
	log.Debug("Client version: %s", hello.ClientVersion)
	welcome := blocky.Handshake(context, hello)
	mustSend(ws, welcome)
	log.Info("Client handshake complete (%s)", welcome.Session)

	// Let the universe deal with the rest.
	ui := universe.GetInteractor(welcome.Session)
	go receiveToInteractor(ws, ui)
	for {
		if packet := ui.Get(); packet != nil {
			log.Debug("Sending packet %T", packet)
			mustSend(ws, packet)
		} else {
			break
		}
	}
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
