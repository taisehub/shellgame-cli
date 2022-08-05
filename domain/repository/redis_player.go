package repository

type PlayerRepository interface {
	GetAll() ([]string, error)
	SetID(string, uint32) error
}
