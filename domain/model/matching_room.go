package model

import (
	"log"
)

var (
	matchingRoom *MatchingRoom = &MatchingRoom{
		MatchingPlayers: make(map[*MatchingPlayer]struct{}),
		message:         make(chan *MatchingMessage),
		register:        make(chan *MatchingPlayer),
		unregister:      make(chan *MatchingPlayer),
	}
)

// shellgame-cliサーバ上で一つだけ存在。
// 対戦待ち状態の管理を行う。
type MatchingRoom struct {
	MatchingPlayers map[*MatchingPlayer]struct{} // 誰がMatchigRoomにいるのか把握するために利用。
	message         chan *MatchingMessage
	register        chan *MatchingPlayer
	unregister      chan *MatchingPlayer
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
	for player := range mr.MatchingPlayers {
		players = append(players, player)
	}
	return players
}

func (mr *MatchingRoom) Run() {
	for {
		select {
		case matchingPlayer := <-mr.register:
			log.Printf("[+] %s entered the room.\n", matchingPlayer.profile.Name)
			mr.MatchingPlayers[matchingPlayer] = struct{}{}
		case matchingPlayer := <-mr.unregister:
			log.Printf("[+] %s exited the room.\n", matchingPlayer.profile.Name)
			if _, ok := mr.MatchingPlayers[matchingPlayer]; ok {
				//matchingPlayerが到達不可能オブジェクトになるはず(願望)
				close(matchingPlayer.matchingChan)
				delete(mr.MatchingPlayers, matchingPlayer)
			}
		case message := <-mr.message:
			log.Printf("[+] %s send message: %s\n", message.source.profile.Name, message.data)
			message.dest.matchingChan <- message
		}
	}
}
