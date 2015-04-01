package main

import (
	"errors"

	"github.com/blixt/geomys"
)

const (
	RegionSize       = 32
	RegionBlockCount = RegionSize * RegionSize
)

type World struct {
	Id       Id
	entities map[string]*Entity
	regions  map[string]*Region
}

func NewWorld() *World {
	return &World{
		Id:       NewId(),
		entities: make(map[string]*Entity),
		regions:  make(map[string]*Region),
	}
}

type Region struct {
	X, Y   int
	Blocks [RegionBlockCount]byte
}

func (w *World) Handler(i *geomys.Interface, event *geomys.Event) error {
	session := i.Context.(*Session)
	if event.Type == "message" {
		switch msg := event.Value.(type) {
		case *MovePlayer:
			return session.Player.Move(msg.X, msg.Y)
		}
	}
	return nil
}

func (w *World) ValidateMove(e *Entity, dx, dy int) error {
	if e.World != w {
		return errors.New("Entity is not of this world")
	}
	// For now, always let the entity move.
	return nil
}
