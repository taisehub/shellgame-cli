package redis

import (
	"context"
	"github.com/taise-hub/shellgame-cli/domain/repository"
	"strconv"
)

type playerRepository struct {
	RedisHandler
}

func NewRedisRepository(rh RedisHandler) repository.PlayerRepository {
	return &playerRepository{rh}
}

func (rep *playerRepository) GetAll() ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return rep.ListGet(ctx, "players")
}

func (rep *playerRepository) SetID(key string, id uint32) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	v := strconv.FormatUint(uint64(id), 10)
	return rep.ListSet(ctx, key, []string{v})
}
