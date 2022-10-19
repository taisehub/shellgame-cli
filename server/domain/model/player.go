package model

import (
	"context"
	"github.com/taise-hub/shellgame-cli/common"
)

type MatchingStatus uint8

const (
	WAITING     MatchingStatus = iota // マッチング待ち状態
	NEGOTIATING                       // マッチング状態(対戦申請を受信または送信中)
)

type MatchingPlayer struct {
	Profile      *common.Profile `json:"profile"`
	Status       MatchingStatus  `json:"status"`
	conn         Conn
	matchingChan chan *common.MatchingMessage
}

func NewMatchingPlayer(id string, name string, conn Conn) *MatchingPlayer {
	profile := &common.Profile{ID: id, Name: name}
	return &MatchingPlayer{
		Profile:      profile,
		Status:       WAITING,
		conn:         conn,
		matchingChan: make(chan *common.MatchingMessage),
	}
}

func (p *MatchingPlayer) GetProfile() *common.Profile {
	return p.Profile
}

func (p *MatchingPlayer) GetStatus() MatchingStatus {
	return p.Status
}

func (p *MatchingPlayer) SetStatus(s MatchingStatus) {
	p.Status = s
}

func (p *MatchingPlayer) GetID() string {
	return p.GetProfile().ID
}

func (p *MatchingPlayer) GetName() string {
	return p.GetProfile().Name
}

func (p *MatchingPlayer) ReadPump() {
	defer func() {
		p.conn.Close()
		GetMatchingRoom().unregister <- p
	}()
	msg := &common.MatchingMessage{}
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
