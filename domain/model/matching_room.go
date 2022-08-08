package model

import (
	"log"
)

var (
	matchingRoom *MatchingRoom = &MatchingRoom{
		MatchingPlayers: make(map[*MatchingPlayer]MatchingState),
		message:         make(chan *MatchingMessage),
		register:        make(chan *MatchingPlayer),
		unregister:      make(chan *MatchingPlayer),
	}
)

// shellgame-cliサーバ上で一つだけ存在。
// 対戦待ち状態の管理を行う。
type MatchingRoom struct {
	MatchingPlayers map[*MatchingPlayer]MatchingState // 誰がMatchigRoomにいるのか把握するために利用。struct{}じゃなくて MatchingStateをvalueにした方がいい？
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
			mr.MatchingPlayers[matchingPlayer] = matchingPlayer.state
		case matchingPlayer := <-mr.unregister:
			log.Printf("[+] %s exited the room.\n", matchingPlayer.profile.Name)
			if _, ok := mr.MatchingPlayers[matchingPlayer]; ok {
				//matchingPlayerが到達不可能オブジェクトになるはず(願望)
				close(matchingPlayer.matchingChan)
				delete(mr.MatchingPlayers, matchingPlayer)
			}
		case message := <-mr.message:
			// messageの中身を確認する。
			// PlayerのState変更は排他処理にする。
			// 1. 申請だった場合
			// 送受信者のMatchingStateが共にWAITINGであれば、受信者に対してリクエストを流す。
			// 送受信者のMatchingStateが共にWAITINGでなければ、Errorを送信者に対してエラーを流す。
			// 2. 申請に対する返答だった場合
			// 申請を受ける場合：送受信者に対して、マッチングしたことを通達する。
			// 申請を断る場合：申請者に対して、マッチングが成立しなかったことを通達する。送受信者のStateをWAITINGにする。
			log.Printf("[+] %s send message: %s\n", message.source.profile.Name, message.data)
			message.dest.matchingChan <- message
		}
	}
}

//TODO: マッチング申請と承諾に関するメソッドはここに書く