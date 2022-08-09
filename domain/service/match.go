package service

import (
	"github.com/taise-hub/shellgame-cli/domain/model"
	"github.com/taise-hub/shellgame-cli/domain/repository"
)

type MatchService struct {
	matchRoomRepo repository.MatchingRoomRepository
}

func NewMatchService(matchRoomRepo repository.MatchingRoomRepository) *MatchService {
	return &MatchService{matchRoomRepo: matchRoomRepo}
}

// MatchingPlayerをMatchingRoomに参加させる
// NOTE: 実装によっては不要になる可能性大
func (srv *MatchService) Wait(player *model.MatchingPlayer) error {
	if err := srv.matchRoomRepo.SetID(player.GetID()); err != nil {
		return err
	}
	return nil
}

// 対戦したい相手に申請を送る
func (srv *MatchService) Request() {
	panic("implement me")
}

// 対戦申請の受け入れ
// playerのStateをPLAYにしてゲームを開始する
func (srv *MatchService) Accept() {
	panic("implement me")
}
