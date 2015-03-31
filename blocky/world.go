package blocky

import (
	"github.com/blixt/geomys"
)

const (
	RegionSize       = 32
	RegionBlockCount = RegionSize * RegionSize
)

type World struct {
	Id          Id
	regionCache map[string]*Region
}

func NewWorld() *World {
	return &World{
		Id:          NewId(),
		regionCache: make(map[string]*Region),
	}
}

type Region struct {
	X, Y   int
	Blocks [RegionBlockCount]byte
}

func (w *World) Handler(i *geomys.Interface, msg interface{}) error {
	// TODO: Implement this handler.
	//session := i.Context.(*Session)
	switch msg := msg.(type) {
	default:
		i.Passthrough(msg)
	}
	return nil
}
