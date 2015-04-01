package main

import (
	"errors"
)

type Entity struct {
	Id    Id
	Kind  string
	World *World `json:"-"`
	X, Y  int
}

func EntityWithId(kind string, id Id) *Entity {
	return &Entity{
		Id:   id,
		Kind: kind,
	}
}

func NewEntity(kind string) *Entity {
	return EntityWithId(kind, NewId())
}

func (e *Entity) GetState() *EntityState {
	return &EntityState{
		Id:   e.Id,
		Kind: e.Kind,
		X:    e.X,
		Y:    e.Y,
	}
}

func (e *Entity) Move(dx, dy int) error {
	if e.World == nil {
		return errors.New("Entity is not in a world")
	}
	if err := e.World.ValidateMove(e, dx, dy); err != nil {
		return err
	}
	e.X += dx
	e.Y += dy
	return nil
}
