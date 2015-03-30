package blocky

type Universe struct {
	activePlayers []*Player
	interactors   []*UniverseInteractor
}

func NewUniverse() *Universe {
	return &Universe{}
}

func (u *Universe) GetInteractor() *UniverseInteractor {
	return NewUniverseInteractor(u)
}

type UniverseInteractor struct {
	universe *Universe
	in       chan interface{}
	out      chan interface{}
}

func NewUniverseInteractor(u *Universe) *UniverseInteractor {
	return &UniverseInteractor{
		universe: u,
		in:       make(chan interface{}),
		out:      make(chan interface{}),
	}
}

func (ui *UniverseInteractor) Close() error {
	ui.out <- &Bye{"Closing connection"}
	close(ui.in)
	close(ui.out)
	return nil
}
