package redis

import (
	"context"
	"github.com/taise-hub/shellgame-cli/server/domain/repository"
)

type matchingRoomRepository struct {
	RedisHandler
}

func NewMatchingRoomRepository(rh RedisHandler) repository.MatchingRoomRepository {
	return &matchingRoomRepository{rh}
}

func (rep *matchingRoomRepository) GetAll() ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return rep.ListGet(ctx, "players")
}

func (rep *matchingRoomRepository) SetID(id string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return rep.ListSet(ctx, "players", []string{id})
}
