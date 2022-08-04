package interfaces


//TODO: /domain/repositoryに移動
type PlayerRepository interface {
}

type playerRepository struct {
	RedisHandler
}

func NewRedisRepository(rh RedisHandler) PlayerRepository {
	return &playerRepository{rh}
}
