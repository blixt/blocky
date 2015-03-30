package blocky

type Universe struct {
	activePlayers []*Player
	interactors   []*UniverseInteractor
}

func NewUniverse() *Universe {
	return &Universe{}
}

func (u *Universe) GetInteractor(s *Session) *UniverseInteractor {
	return NewUniverseInteractor(u, s)
}

type UniverseInteractor struct {
	universe *Universe
	session  *Session
	in       chan interface{}
	out      chan interface{}
}

func NewUniverseInteractor(u *Universe, s *Session) *UniverseInteractor {
	return &UniverseInteractor{
		universe: u,
		session:  s,
		in:       make(chan interface{}),
		out:      make(chan interface{}),
	}
}

func (ui *UniverseInteractor) Close() {
	close(ui.in)
	close(ui.out)
}

func (ui *UniverseInteractor) Get() interface{} {
	return <-ui.out
}

func (ui *UniverseInteractor) Put(packet interface{}) {
	ui.in <- packet
}
