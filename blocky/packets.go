package blocky

type Bye struct {
	Message string
}

type Hello struct {
	SessionId     Id
	ClientVersion string
}

type Ping struct {
	Time int64
}

type Welcome struct {
	Session       *Session
	ServerVersion string
}

func Handshake(hello *Hello) *Welcome {
	return &Welcome{
		Session:       GetOrCreateSession(hello.SessionId),
		ServerVersion: "0.1.0.001",
	}
}
