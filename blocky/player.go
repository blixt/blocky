package blocky

type Player struct {
	Id   Id
	Name string
}

var players = make(map[Id]*Player)

func GetPlayer(id Id) *Player {
	return players[id]
}

func NewPlayer() *Player {
	player := &Player{Id: NewId(), Name: "Guest"}
	players[player.Id] = player
	return player
}
