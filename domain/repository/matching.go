package repository

// 対戦待ちのプレイヤーに関する操作を行うRepository
type MatchingRepository interface {
	GetAll() ([]string, error)
	SetID(uint32) error
}
