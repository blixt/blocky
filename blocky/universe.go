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

func (u *Universe) PutAll(packet interface{}) {
	count, deleted := len(u.interactors), 0
	for i, ui := range u.interactors {
		if err := ui.Put(packet); err != nil {
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
		u.PutAll(&Ping{time.Now().Unix() * 1000})
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

func (ui *UniverseInteractor) Get() interface{} {
	return <-ui.channel
}

func (ui *UniverseInteractor) Put(packet interface{}) error {
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
