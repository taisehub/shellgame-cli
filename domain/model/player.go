package model

import (
	"time"
	"context"
)

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

func (p *MatchingPlayer) WritePump() {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)
	defer func() {
		cancel()
		p.conn.Close()
	}()
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
