package repository

// 対戦待ちのプレイヤーに関する操作を行うRepository
type MatchingRoomRepository interface {
	GetAll() ([]string, error)
	SetID(string) error
}
