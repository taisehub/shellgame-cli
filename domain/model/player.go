package model

type Player interface {
	GetProfile() *Profile
}

type MatchingPlayer struct {
	profile      *Profile
	state        State
	conn         Conn
	matchingChan chan *MatchingMessage
}

func NewMatchingPlayer(profile *Profile) *MatchingPlayer {
	return &MatchingPlayer{
		profile:      profile,
		state:        INACTIVE,
		matchingChan: make(chan *MatchingMessage),
	}
}

// FIXME: IDではなくProfileを返すようにする
func (p *MatchingPlayer) GetProfile() *Profile {
	return p.profile
}

// TODO: context対応
func (p *MatchingPlayer) ReadPump() {
	defer func() {
		p.conn.Close()
		GetMatchingRoom().unregister <- p
	}()
	var msg  *MatchingMessage
	for {
		if err := p.conn.Read(msg); err != nil {
			return
		}
		GetMatchingRoom().message <- msg
	}
}

// TODO: context対応
func (p *MatchingPlayer) WritePump() {
	defer p.conn.Close()
	for {
		msg, ok := <- p.matchingChan
		if !ok {
			return
		}
		if err := p.conn.Write(msg); err != nil {
			return
		}
	}
}
