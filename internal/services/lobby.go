package services

import (
	"github.com/louisfield/multi_backend/internal/types"
	"github.com/lxzan/gws"
)

var lobbies = []*types.Lobby{}

func MaybeCreateLobby(session gws.SessionStorage, name string) {
	lobby := FindLobby(name)
	if lobby != nil {
		return
	}
	new_lobby := newLobby(name)

	addLobby(new_lobby)
}

func addLobby(lobby *types.Lobby) {
	lobbies = append(lobbies, lobby)
}

func newLobby(name string) *types.Lobby {
	return &types.Lobby{
		Token: name,
		Users: []*types.User{},
	}
}

func FindLobby(token string) *types.Lobby {
	for _, lobby := range lobbies {
		if lobby.Token == token {
			return lobby
		}
	}
	return nil
}
