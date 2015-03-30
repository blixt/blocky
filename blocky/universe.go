package blocky

import (
	"errors"
	"time"
)

type Universe struct {
	auth       *Authenticator
	interfaces []*Interface
}

func NewUniverse() *Universe {
	return &Universe{}
}

func (u *Universe) Handle(i *Interface, packet interface{}) error {
	return errors.New("Not implemented")
}

func (u *Universe) NewInterface(c Context) *Interface {
	i := NewInterface(c, u)
	i.PushHandler(u.auth)
	u.interfaces = append(u.interfaces, i)
	return i
}

func (u *Universe) Run() {
	for {
		u.putAll(&Ping{time.Now().Unix() * 1000})
		time.Sleep(5 * time.Second)
	}
}

func (u *Universe) putAll(packet interface{}) {
	count, deleted := len(u.interfaces), 0
	for index, i := range u.interfaces {
		if err := i.putClient(packet); err != nil {
			// Forget this interface because it's not active anymore.
			deleted++
			u.interfaces[index] = u.interfaces[count-deleted]
		}
	}
	if deleted > 0 {
		// Ensure that we don't keep garbage references around.
		for index := deleted; index > 0; index-- {
			u.interfaces[count-index] = nil
		}
		// Shorten the slice.
		u.interfaces = u.interfaces[:count-deleted]
	}
}
