package blocky

import (
	"errors"
)

type Handler interface {
	Handle(i *Interface, packet interface{}) error
}

type Interface struct {
	Context  Context
	session  *Session
	handlers []Handler
	open     bool
	channel  chan interface{}
}

func NewInterface(c Context, h Handler) *Interface {
	return &Interface{
		Context:  c,
		handlers: []Handler{h},
		open:     true,
		channel:  make(chan interface{}, 10),
	}
}

func (i *Interface) Close() {
	i.handlers = nil
	i.open = false
	close(i.channel)
}

// Gets a packet for the client (or waits until one is available).
func (i *Interface) Get() interface{} {
	return <-i.channel
}

func (i *Interface) PopHandler() {
	if len(i.handlers) < 2 {
		panic("Cannot pop root handler")
	}
	i.handlers[len(i.handlers)-1] = nil
	i.handlers = i.handlers[:len(i.handlers)-1]
}

func (i *Interface) PushHandler(h Handler) {
	i.handlers = append(i.handlers, h)
}

// Handles a packet from the client.
func (i *Interface) Put(packet interface{}) error {
	if !i.open {
		return errors.New("The interface is closed")
	}
	return i.handlers[len(i.handlers)-1].Handle(i, packet)
}

// Sends a packet to the client.
func (i *Interface) putClient(packet interface{}) error {
	if !i.open {
		return errors.New("The interface is closed")
	}
	select {
	case i.channel <- packet:
		return nil
	default:
		i.Close()
		return errors.New("Interface overflowed with packets")
	}
}
