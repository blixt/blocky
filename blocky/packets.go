package blocky

type Hello struct {
	SessionId     Id
	ClientVersion string
}

type Welcome struct {
	Session       *Session
	ServerVersion string
}
