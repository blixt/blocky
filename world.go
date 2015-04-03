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
	world := &World{
		Id:       NewId(),
		entities: make(map[Id]*Entity),
		regions:  make(map[string]*Region),
	}
	// TODO: Create regions dynamically.
	region := NewRegion(0, 0)
	for i := 0; i < RegionSize; i++ {
		region.SetBlock(i, 0, 1)
		region.SetBlock(0, i, 1)
		region.SetBlock(i, RegionSize-1, 1)
		region.SetBlock(RegionSize-1, i, 1)
	}
	world.AddRegion(region)
	return world
}

type Region struct {
	X, Y   int
	Blocks [RegionBlockCount]byte
}

func NewRegion(x, y int) *Region {
	region := new(Region)
	region.X = x
	region.Y = y
	return region
}

func RegionKey(x, y int) string {
	return fmt.Sprintf("%d:%d", x, y)
}

func (r *Region) Key() string {
	return RegionKey(r.X, r.Y)
}

func (r *Region) SetBlock(x, y int, block byte) {
	if x < 0 || x >= RegionSize || y < 0 || y >= RegionSize {
		panic("Region coordinates out of bounds")
	}
	r.Blocks[y*RegionSize+x] = block
}

func (w *World) AddEntity(e *Entity) error {
	if _, ok := w.entities[e.Id]; ok {
		return errors.New("That entity is already in this world")
	}
	w.entities[e.Id] = e
	return nil
}

func (w *World) AddRegion(r *Region) error {
	key := r.Key()
	if _, ok := w.regions[key]; ok {
		return errors.New("That region space is already occupied in this world")
	}
	w.regions[key] = r
	return nil
}

func (w *World) GetRegion(x, y int) (region *Region, err error) {
	region = w.regions[RegionKey(x, y)]
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
