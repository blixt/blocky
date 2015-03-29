package blocky

type Hello struct {
	SessionId     Id
	ClientVersion string
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
