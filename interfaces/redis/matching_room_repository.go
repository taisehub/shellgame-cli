package redis

import (
	"context"
	"github.com/taise-hub/shellgame-cli/domain/repository"
	"strconv"
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

func (rep *matchingRoomRepository) SetID(id uint32) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	v := strconv.FormatUint(uint64(id), 10)
	return rep.ListSet(ctx, "players", []string{v})
}
