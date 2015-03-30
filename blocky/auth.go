package blocky

import (
	"errors"
)

type Authenticator struct {
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{}
}

func (a *Authenticator) Handle(i *Interface, packet interface{}) error {
	if hello, ok := packet.(*Hello); ok {
		// Shake hands.
		welcome := Handshake(i.Context, hello)
		i.putClient(welcome)
		i.PopHandler()
		return nil
	} else {
		return errors.New("Expected a Hello packet")
	}
}
