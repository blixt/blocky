package main

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
