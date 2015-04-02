package main

import (
	"fmt"
	"time"

	"github.com/blixt/geomys"
	"golang.org/x/net/websocket"
)

// TODO: The universe should support many worlds.
type Universe struct {
	geomys.WebSocketServerBase
	Server   *geomys.Server
	World    *World
	lastPing *Ping
}

func NewUniverse() *Universe {
	return &Universe{
		Server: geomys.NewServer(),
		World:  NewWorld(),
	}
}

func (u *Universe) GetInterface(ws *websocket.Conn) *geomys.Interface {
	i := u.Server.NewInterface(nil)
	i.PushHandler(u.Handler)
	i.PushHandler(u.World.Handler)
	i.PushHandler(SessionHandler)
	return i
}

func (u *Universe) GetMessage(msgType string) (interface{}, error) {
	switch msgType {
	case "Bye":
		return new(Bye), nil
	case "Hello":
		return new(Hello), nil
	case "LoadRegion":
		return new(LoadRegion), nil
	case "Ping":
		return new(Ping), nil
	default:
		return nil, fmt.Errorf("Unsupported message type %s", msgType)
	}
}

func (u *Universe) Handler(i *geomys.Interface, event *geomys.Event) error {
	session := i.Context.(*Session)
	switch event.Type {
	case "message":
		switch msg := event.Value.(type) {
		case *Ping:
			if msg.Id == u.lastPing.Id {
				session.Player.latency = NewPing().Time - u.lastPing.Time
				session.Player.lastPing = time.Now()
			}
		default:
			return fmt.Errorf("Unexpected message %T", msg)
		}
	}
	return nil
}

func (u *Universe) Run() {
	for {
		u.lastPing = NewPing()
		u.Server.SendAll(u.lastPing)
		time.Sleep(5 * time.Second)
	}
}
