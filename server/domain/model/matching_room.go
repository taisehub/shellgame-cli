package model

import (
	"fmt"
	"github.com/taise-hub/shellgame-cli/common"
	"log"
	"sync"
)

var (
	mu sync.Mutex
	matchingRoom *MatchingRoom = &MatchingRoom{
		Players:    make(map[string]*MatchingPlayer),
		message:    make(chan *common.MatchingMessage),
		register:   make(chan *MatchingPlayer),
		unregister: make(chan *MatchingPlayer),
	}
)

// shellgame-cliサーバ上で一つだけ存在。
// 対戦待ち状態の管理を行う。
type MatchingRoom struct {
	Players    map[string]*MatchingPlayer // 誰がMatchigRoomにいるのか把握するために利用。
	message    chan *common.MatchingMessage
	register   chan *MatchingPlayer
	unregister chan *MatchingPlayer
}

func GetMatchingRoom() *MatchingRoom {
	return matchingRoom
}

func (mr *MatchingRoom) GetRegisterChan() chan<- *MatchingPlayer {
	return mr.register
}

func (mr *MatchingRoom) GetUnregisterChan() chan<- *MatchingPlayer {
	return mr.unregister
}

func (mr *MatchingRoom) GetMatchingPlayers() []*MatchingPlayer {
	var players []*MatchingPlayer
	for _, player := range mr.Players {
		players = append(players, player)
	}
	return players
}

func (mr *MatchingRoom) Run() {
	for {
		select {
		case player := <-mr.register:
			log.Printf("[+] %s entered the room.\n", player.GetName())
			for _, p := range mr.Players {
				var msg = &common.MatchingMessage{
					Source: player.GetProfile(),
					Dest:   nil,
					Data:   common.JOIN,
				}
				// 参加は全員に送信する
				p.matchingChan <- msg
			}
			mr.enterRoom(player)
		case player := <-mr.unregister:
			mr.exitRoom(player)
			for _, p := range mr.Players {
				var msg = &common.MatchingMessage{
					Source: player.GetProfile(),
					Dest:   nil,
					Data:   common.LEAVE,
				}
				// 退室は全員に送信する
				p.matchingChan <- msg
			}
		case msg := <-mr.message:
			switch msg.Data {
			case common.OFFER:
				log.Printf("[+] OFFER: %s to %s\n", mr.Players[msg.Source.ID].GetName(), mr.Players[msg.Dest.ID].GetName())
				mr.HandleOffer(msg)
				mr.Players[msg.Dest.ID].matchingChan <- msg
			case common.CANCEL_OFFER:
				log.Printf("[+] CANCEL OFFER: %s to %s\n", mr.Players[msg.Source.ID].GetName(), mr.Players[msg.Dest.ID].GetName())
				mr.HandleCancelOffer(msg)
				mr.Players[msg.Dest.ID].matchingChan <- msg
			case common.ACCEPT:
				log.Printf("[+] ACCEPT OFFER: %s to %s\n", mr.Players[msg.Source.ID].GetName(), mr.Players[msg.Dest.ID].GetName())
				mr.HandleAccept(msg)
				mr.Players[msg.Source.ID].matchingChan <- msg
				mr.Players[msg.Dest.ID].matchingChan <- msg
			case common.DENY:
				log.Printf("[+] DENY OFFER: %s to %s\n", mr.Players[msg.Source.ID].GetName(), mr.Players[msg.Dest.ID].GetName())
				mr.HandleDeny(msg)
				mr.Players[msg.Dest.ID].matchingChan <- msg
			default:
				err := &common.MatchingMessage{Data: common.ERROR}
				mr.Players[msg.Source.ID].matchingChan <- err
			}
		}
	}
}

func (mr *MatchingRoom) enterRoom(p *MatchingPlayer) {
	mr.Players[p.GetID()] = p
}

func (mr *MatchingRoom) exitRoom(p *MatchingPlayer) {
	log.Printf("[+] %s exited the room.\n", p.GetName())
	if _, ok := mr.Players[p.GetID()]; ok {
		close(p.matchingChan)
		delete(mr.Players, p.GetID())
	}
}

// 対戦申請処理
func (mr *MatchingRoom) HandleOffer(msg *common.MatchingMessage) {
	// 受信者がRoomにいることを確認
	_, ok := mr.Players[msg.Dest.ID]
	if !ok {
		err := &common.MatchingMessage{Data: common.ERROR}
		mr.Players[msg.Source.ID].matchingChan <- err
		return
	}
	mr.changeToNegotiating(msg.Source, msg.Dest)
}

// 対戦申請キャンセル処理
func (mr *MatchingRoom) HandleCancelOffer(msg *common.MatchingMessage) {
	_, ok := mr.Players[msg.Dest.ID]
	if !ok {
		err := &common.MatchingMessage{Data: common.ERROR}
		mr.Players[msg.Source.ID].matchingChan <- err
		return
	}

	mr.changeToWaiting(msg.Source, msg.Dest)
}

// 対戦申請に対する承諾処理
func (mr *MatchingRoom) HandleAccept(msg *common.MatchingMessage) {
	_, ok := mr.Players[msg.Dest.ID]
	if !ok {
		err := &common.MatchingMessage{Data: common.ERROR}
		mr.Players[msg.Source.ID].matchingChan <- err
		return
	}

	isNego := func(src, dst *common.Profile) error {
		mu.Lock()
		defer mu.Unlock()
		if mr.Players[src.ID].GetStatus() == WAITING {
			return fmt.Errorf("souce player is not NEGOTIATING")
		} else if mr.Players[dst.ID].GetStatus() == WAITING {
			return fmt.Errorf("destination player is not NEGOTIATING")
		}
		return nil
	}

	if err := isNego(msg.Source, msg.Dest); err != nil {
		err := &common.MatchingMessage{Data: common.ERROR}
		mr.Players[msg.Source.ID].matchingChan <- err
	}
}

// 対戦申請に対する不承諾処理
func (mr *MatchingRoom) HandleDeny(msg *common.MatchingMessage) {
	_, ok := mr.Players[msg.Dest.ID]
	if !ok {
		err := &common.MatchingMessage{Data: common.ERROR}
		mr.Players[msg.Source.ID].matchingChan <- err
		return
	}

	mr.changeToWaiting(msg.Source, msg.Dest)
}

// 申請者(src)と承諾者(dst)のステータスが共にWAITINGであること確認し、両者のステータスをNEGOTIATINGにする。
func (mr *MatchingRoom) changeToNegotiating(src, dst *common.Profile) error {
	if mr.Players[src.ID] == nil {
		return fmt.Errorf("source player is not in the room.")
	} else if mr.Players[dst.ID] == nil {
		return fmt.Errorf("destination player is not in the room.")
	}

	mu.Lock()
	defer mu.Unlock()

	if mr.Players[src.ID].GetStatus() == NEGOTIATING {
		return fmt.Errorf("souce player has already been NEGOTIATING")
	} else if mr.Players[dst.ID].GetStatus() == NEGOTIATING {
		return fmt.Errorf("destination player has alredy been NEGOTIATING")
	}

	mr.Players[src.ID].SetStatus(NEGOTIATING)
	mr.Players[dst.ID].SetStatus(NEGOTIATING)
	return nil
}

// 申請者(src)と承諾者(dst)のステータスが共にNEGITIATINGであること確認し、両者のステータスをWAITINGにする。
func (mr *MatchingRoom) changeToWaiting(src, dst *common.Profile) error {
	if mr.Players[src.ID] == nil {
		return fmt.Errorf("source player is not in the room.")
	} else if mr.Players[dst.ID] == nil {
		return fmt.Errorf("destination player is not in the room.")
	}

	mu.Lock()
	defer mu.Unlock()

	// FIXME: 交渉中でないPlayerを宛先にしてMessageを送信することで交渉解除することが可能。
	if mr.Players[src.ID].GetStatus() == WAITING {
		return fmt.Errorf("souce player has already been WAITING")
	} else if mr.Players[dst.ID].GetStatus() == WAITING {
		return fmt.Errorf("destination player has already been WAITING")
	}

	mr.Players[src.ID].SetStatus(WAITING)
	mr.Players[dst.ID].SetStatus(WAITING)
	return nil
}