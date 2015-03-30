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

func Handshake(context Context, hello *Hello) *Welcome {
	return &Welcome{
		Session:       GetOrCreateSession(hello.SessionId),
		ServerVersion: context.Version,
	}
}
