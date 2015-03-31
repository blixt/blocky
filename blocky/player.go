package blocky

import (
	"time"
)

type Player struct {
	Id       Id
	Name     string
	Ping     float64
	PingTime time.Time
}

var players = make(map[Id]*Player)

func GetPlayer(id Id) *Player {
	return players[id]
}

func NewPlayer() *Player {
	player := &Player{Id: NewId(), Name: "Guest"}
	players[player.Id] = player
	return player
}
