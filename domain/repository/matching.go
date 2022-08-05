package repository

type MatchingRepository interface {
	GetAll() ([]string, error)
	SetID(string, uint32) error
}
