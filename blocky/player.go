package blocky

import (
	"errors"
	"time"
)

type Player struct {
	Id       Id
	Name     string
	Ping     float64
	PingTime time.Time
	X, Y     int
	lastMove time.Time
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

func (p *Player) Move(dx, dy int) error {
	if dx < -1 || dx > 1 || dy < -1 || dy > 1 {
		return errors.New("Invalid move delta")
	} else if dx != 0 && dy != 0 {
		return errors.New("Can only move in one direction at a time")
	} else if time.Now().Sub(p.lastMove) < 100*time.Millisecond {
		return errors.New("Can only move once every 100 ms")
	}
	p.X += dx
	p.Y += dy
	p.lastMove = time.Now()
	return nil
}
