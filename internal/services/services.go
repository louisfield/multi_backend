package services

import (
	"github.com/louisfield/multi_backend/internal/types"
)

func NewUser(name string) *types.User {
	return &types.User{
		Name: name,
	}
}
