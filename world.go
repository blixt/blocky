package main

import (
	"errors"
	"fmt"

	"github.com/blixt/geomys"
)

const (
	RegionSize       = 32
	RegionBlockCount = RegionSize * RegionSize
)

type World struct {
	Id             Id
	StartX, StartY int
	entities       map[Id]*Entity
	regions        map[string]*Region
}

func NewWorld() *World {
	return &World{
		Id:       NewId(),
		entities: make(map[Id]*Entity),
		regions:  make(map[string]*Region),
	}
}

type Region struct {
	X, Y   int
	Blocks [RegionBlockCount]byte
}

func (w *World) AddEntity(e *Entity) error {
	if _, ok := w.entities[e.Id]; ok {
		return errors.New("That entity is already in this world")
	}
	w.entities[e.Id] = e
	return nil
}

func (w *World) GetRegion(x, y int) (region *Region, err error) {
	region = w.regions[fmt.Sprintf("%d:%d", x, y)]
	if region == nil {
		err = errors.New("The requested region does not exist")
	}
	return
}

func (w *World) Handler(i *geomys.Interface, event *geomys.Event) error {
	session := i.Context.(*Session)
	switch event.Type {
	case "auth":
		// Add the player to the world once they're authed.
		session.Player.GoToWorld(w)
		i.Send(NewEnterWorld(session.Player))
	case "message":
		switch msg := event.Value.(type) {
		case *LoadRegion:
			if msg.WorldId != w.Id {
				return errors.New("Player sent the wrong world id")
			}
			if region, err := w.GetRegion(msg.X, msg.Y); err != nil {
				return err
			} else {
				i.Send(region)
			}
		case *MovePlayer:
			return session.Player.Move(msg.X, msg.Y)
		}
	}
	return nil
}

func (w *World) RemoveEntity(e *Entity) error {
	if _, ok := w.entities[e.Id]; !ok {
		return errors.New("That entity is not in this world")
	}
	delete(w.entities, e.Id)
	return nil
}

func (w *World) ValidateMove(e *Entity, dx, dy int) error {
	if e.World != w {
		return errors.New("Entity is not of this world")
	}
	// For now, always let the entity move.
	return nil
}
