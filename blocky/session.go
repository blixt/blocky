package blocky

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
