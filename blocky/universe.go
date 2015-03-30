package blocky

import (
	"errors"
	"time"
)

type Universe struct {
	activePlayers []*Player
	interactors   []*UniverseInteractor
}

func NewUniverse() *Universe {
	return &Universe{}
}

func (u *Universe) GetInteractor(s *Session) *UniverseInteractor {
	ui := NewUniverseInteractor(u, s)
	u.interactors = append(u.interactors, ui)
	return ui
}

func (u *Universe) Handle(ui *UniverseInteractor, packet interface{}) {
}

func (u *Universe) putAll(packet interface{}) {
	count, deleted := len(u.interactors), 0
	for i, ui := range u.interactors {
		if err := ui.put(packet); err != nil {
			// Forget this interactor because it's not active anymore.
			deleted++
			u.interactors[i] = u.interactors[count-deleted]
		}
	}
	if deleted > 0 {
		// Ensure that we don't keep garbage references around.
		for i := deleted; i > 0; i-- {
			u.interactors[count-i] = nil
		}
		// Shorten the slice.
		u.interactors = u.interactors[:count-deleted]
	}
}

func (u *Universe) Run() {
	for {
		u.putAll(&Ping{time.Now().Unix() * 1000})
		time.Sleep(5 * time.Second)
	}
}

type UniverseInteractor struct {
	universe *Universe
	session  *Session
	open     bool
	channel  chan interface{}
}

func NewUniverseInteractor(u *Universe, s *Session) *UniverseInteractor {
	return &UniverseInteractor{
		universe: u,
		session:  s,
		open:     true,
		channel:  make(chan interface{}, 10),
	}
}

func (ui *UniverseInteractor) Close() {
	ui.open = false
	close(ui.channel)
}

// Gets a packet for the client (or waits until one is available).
func (ui *UniverseInteractor) Get() interface{} {
	return <-ui.channel
}

// Handles a packet from the client.
func (ui *UniverseInteractor) Handle(packet interface{}) error {
	if !ui.open {
		return errors.New("The universe interactor is closed")
	}
	ui.universe.Handle(ui, packet)
	return nil
}

// Sends a packet to the client.
func (ui *UniverseInteractor) put(packet interface{}) error {
	if !ui.open {
		return errors.New("The universe interactor is closed")
	}
	select {
	case ui.channel <- packet:
		return nil
	default:
		ui.Close()
		return errors.New("Universe interactor overflowed with packets")
	}
}
