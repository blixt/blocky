package main

import (
	"time"
)

// Tells the client/server that the connection is about to be closed.
type Bye struct {
	Message string
}

// Lets the client know that the player entered a world.
type EnterWorld struct {
	Player   *Player
	World    *World
	Entities []*Entity
}

func NewEnterWorld(p *Player) *EnterWorld {
	entities := make([]*Entity, 0)
	for id, entity := range p.World.entities {
		// TODO: Filter so that only nearby entities are announced.
		if id != p.Id {
			entities = append(entities, entity)
		}
	}
	return &EnterWorld{p, p.World, entities}
}

// Notifies the client of the state of a new or existing entity.
type EntityState struct {
	Id    Id
	Kind  string
	X, Y  int
	State int
}

// An identification packet from the client.
type Hello struct {
	SessionId     Id
	ClientVersion string
}

// A request from the client to load a region.
type LoadRegion struct {
	WorldId Id
	X, Y    int
}

// Requests that the player be moved.
type MovePlayer struct {
	X, Y int
}

// Determines the latency between the client and the server.
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

// Greets the client once it's identified itself.
type Welcome struct {
	Session       *Session
	ServerVersion string
}
