package types

import (
	"github.com/lxzan/gws"
)

type User struct {
	Name string
	Conn *gws.Conn
}

type Input struct {
	Message string `json:"message"`
	Event   string `json:"event"`
}

type Lobby struct {
	Users          []*User
	MaxPlayers     int
	CurrentPlayers int
	Token          string
}

func (l *Lobby) Join(user *User) {
	l.Users = append(l.Users, user)
	l.CurrentPlayers++
}

func (l *Lobby) Broadcast(message string) {
	for _, user := range l.Users {
		go user.Conn.WriteString(message)
	}
}
