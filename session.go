package main

import (
	"fmt"

	"github.com/blixt/geomys"
)

type Session struct {
	Id     Id
	Player *Player
}

var sessions = make(map[Id]*Session)

func GetOrCreateSession(id Id) *Session {
	if session := GetSession(id); session != nil {
		return session
	} else {
		return NewSession()
	}
}

func GetSession(id Id) *Session {
	return sessions[id]
}

func NewSession() *Session {
	session := &Session{Id: NewId(), Player: NewPlayer()}
	sessions[session.Id] = session
	return session
}

func (s *Session) String() string {
	return fmt.Sprintf("Session: %s Player: %s Name: %s", s.Id, s.Player.Id, s.Player.Name)
}

func Handshake(hello *Hello) (*Welcome, error) {
	welcome := &Welcome{
		Session:       GetOrCreateSession(hello.SessionId),
		ServerVersion: "0.1.0.001",
	}
	// TODO: Return an error if handshake fails.
	return welcome, nil
}

// Handles incoming "Hello" messages, discards everything else.
func SessionHandler(i *geomys.Interface, event *geomys.Event) error {
	if event.Type == "message" {
		switch msg := event.Value.(type) {
		case *Hello:
			// Shake hands.
			if welcome, err := Handshake(msg); err != nil {
				return err
			} else {
				// Handshake succeeded, set session context and let client know.
				i.Context = welcome.Session
				i.Send(welcome)
				// Remove this handler now that the user has been authenticated.
				i.RemoveHandler()
				event.StopPropagation()
				// Let other handlers know that the user has been authenticated.
				i.Dispatch(geomys.NewEvent("auth", welcome.Session))
			}
		default:
			return fmt.Errorf("Unexpected message %T", msg)
		}
	}
	return nil
}
