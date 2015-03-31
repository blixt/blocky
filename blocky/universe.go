package blocky

import (
	"fmt"
	"time"

	"github.com/blixt/geomys"
	"golang.org/x/net/websocket"
)

type Universe struct {
	geomys.WebSocketServerBase
	Server   *geomys.Server
	lastPing *Ping
}

func NewUniverse() *Universe {
	return &Universe{
		Server: geomys.NewServer(),
	}
}

func (u *Universe) GetInterface(ws *websocket.Conn) *geomys.Interface {
	i := u.Server.NewInterface(nil)
	i.PushHandler(u.handler)
	i.PushHandler(SessionHandler)
	return i
}

func (u *Universe) GetMessage(msgType string) (interface{}, error) {
	switch msgType {
	case "Hello":
		return new(Hello), nil
	case "Ping":
		return new(Ping), nil
	default:
		return nil, fmt.Errorf("Unsupported message type %s", msgType)
	}
}

func (u *Universe) Run() {
	for {
		u.lastPing = NewPing()
		u.Server.SendAll(u.lastPing)
		time.Sleep(5 * time.Second)
	}
}

func (u *Universe) handler(i *geomys.Interface, msg interface{}) error {
	session := i.Context.(*Session)
	switch msg := msg.(type) {
	case *Ping:
		if msg.Id == u.lastPing.Id {
			session.Player.Ping = NewPing().Time - u.lastPing.Time
			session.Player.PingTime = time.Now()
		}
	default:
		return fmt.Errorf("Unexpected message %T", msg)
	}
	return nil
}
