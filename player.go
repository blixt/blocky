package main

import (
	"errors"
	"time"
)

type Player struct {
	*Entity
	Name     string
	Ping     float64   `json:"-"`
	PingTime time.Time `json:"-"`
	lastMove time.Time
}

var players = make(map[Id]*Player)

func GetPlayer(id Id) *Player {
	return players[id]
}

func NewPlayer() *Player {
	player := &Player{Entity: NewEntity("player"), Name: "Guest"}
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
	if err := p.Entity.Move(dx, dy); err != nil {
		return err
	}
	p.lastMove = time.Now()
	return nil
}
