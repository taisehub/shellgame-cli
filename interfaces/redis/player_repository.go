package redis

import (
	"context"
	"github.com/taise-hub/shellgame-cli/domain/repository"
	"strconv"
)

type matchingRepository struct {
	RedisHandler
}

func NewmatchingRepository(rh RedisHandler) repository.MatchingRepository {
	return &matchingRepository{rh}
}

func (rep *matchingRepository) GetAll() ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return rep.ListGet(ctx, "players")
}

func (rep *matchingRepository) SetID(key string, id uint32) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	v := strconv.FormatUint(uint64(id), 10)
	return rep.ListSet(ctx, key, []string{v})
}
