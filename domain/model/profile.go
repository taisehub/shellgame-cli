package model

type Profile struct {
	id   uint32
	name string
}

func (p *Profile) GetID() uint32 {
	return p.id
}
