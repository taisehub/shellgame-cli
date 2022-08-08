package model

import (
	"context"
)

type Player interface {
	GetProfile() *Profile
}

type MatchingPlayer struct {
	profile      *Profile
	state        MatchingState
	conn         Conn
	matchingChan chan *MatchingMessage
}

func NewMatchingPlayer(id uint32, name string, conn Conn) *MatchingPlayer {
	profile := &Profile{ID: id, Name: name}
	return &MatchingPlayer{
		profile:      profile,
		state:        WAITING,
		conn:         conn,
		matchingChan: make(chan *MatchingMessage),
	}
}

func (p *MatchingPlayer) GetProfile() *Profile {
	return p.profile
}

func (p *MatchingPlayer) ReadPump() {
	defer func() {
		p.conn.Close()
		GetMatchingRoom().unregister <- p
	}()
	var msg *MatchingMessage
	for {
		if err := p.conn.Read(msg); err != nil {
			return
		}
		GetMatchingRoom().message <- msg
	}
}

func (p *MatchingPlayer) WritePump(ctx context.Context) {
	defer p.conn.Close()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-p.matchingChan:
			if !ok {
				return
			}
			if err := p.conn.Write(msg); err != nil {
				return
			}
		}
	}
}
