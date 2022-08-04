package model

type Profile struct {
	id   uint32
	name string
}

func NewProfile(id uint32, name string) *Profile {
	return &Profile{
		id:   id,
		name: name,
	}
}
