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

func (p *Player) GetID() uint32 {
	return p.profile.GetID()
}