package model

import (
	"context"
)

type MatchingStatus uint8

const (
	WAITING     MatchingStatus = iota // マッチング待ち状態
	NEGOTIATING                       // マッチング状態(対戦申請を受信または送信中)
)

type Player interface {
	GetID() uint32
}

type MatchingPlayer struct {
	Profile      *Profile              `json:"profile"`
	Status       MatchingStatus        `json:"status"`
	conn         Conn                  `json:"-"`
	matchingChan chan *MatchingMessage `json:"-"`
}

func NewMatchingPlayer(id uint32, name string, conn Conn) *MatchingPlayer {
	profile := &Profile{ID: id, Name: name}
	return &MatchingPlayer{
		Profile:      profile,
		Status:       WAITING,
		conn:         conn,
		matchingChan: make(chan *MatchingMessage),
	}
}

func (p *MatchingPlayer) GetID() uint32 {
	return p.Profile.ID
}

func (p *MatchingPlayer) GetName() string {
	return p.Profile.Name
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
