package blocky

import (
	"fmt"
	"time"

	"github.com/blixt/geomys"
	"golang.org/x/net/websocket"
)

type Universe struct {
	geomys.WebSocketServerBase
	Server *geomys.Server
}

func NewUniverse() *Universe {
	return &Universe{
		Server: geomys.NewServer(),
	}
}

func (u *Universe) GetInterface(ws *websocket.Conn) *geomys.Interface {
	i := u.Server.NewInterface(nil)
	i.PushHandler(u.handleDefault)
	i.PushHandler(u.handleAuth)
	return i
}

func (u *Universe) GetMessage(msgType string) (interface{}, error) {
	switch msgType {
	case "Hello":
		return new(Hello), nil
	default:
		return nil, fmt.Errorf("Unsupported message type %s", msgType)
	}
}

func (u *Universe) Run() {
	for {
		u.Server.SendAll(&Ping{time.Now().Unix() * 1000})
		time.Sleep(5 * time.Second)
	}
}

func (u *Universe) handleAuth(i *geomys.Interface, msg interface{}) error {
	if hello, ok := msg.(*Hello); ok {
		// Shake hands.
		welcome := Handshake(hello)
		i.Send(welcome)
		i.PopHandler()
		return nil
	} else {
		return fmt.Errorf("Expected a Hello message, got %T", msg)
	}
}

func (u *Universe) handleDefault(i *geomys.Interface, msg interface{}) error {
	return nil
}
