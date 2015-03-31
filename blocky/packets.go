package blocky

import (
	"time"
)

type Bye struct {
	Message string
}

type Hello struct {
	SessionId     Id
	ClientVersion string
}

type Ping struct {
	Id   Id
	Time float64
}

func NewPing() *Ping {
	return &Ping{
		NewId(),
		float64(time.Now().UnixNano()) / float64(time.Millisecond),
	}
}

type Welcome struct {
	Session       *Session
	ServerVersion string
}

func Handshake(hello *Hello) *Welcome {
	return &Welcome{
		Session:       GetOrCreateSession(hello.SessionId),
		ServerVersion: "0.1.0.001",
	}
}
