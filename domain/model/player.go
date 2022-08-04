package model

type Player struct {
	profile Profile
	state   State
}

func NewPlayer(profile Profile) *Player {
	return &Player{
		profile: profile,
		state:   INACTIVE,
	}
}
