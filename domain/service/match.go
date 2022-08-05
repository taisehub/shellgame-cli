package service

import (
	"github.com/taise-hub/shellgame-cli/domain/model"
	"github.com/taise-hub/shellgame-cli/domain/repository"
)

type MatchService struct {
	matchRepo repository.MatchingRepository
}

func NewMatchService(matchRepo repository.MatchingRepository) *MatchService {
	return &MatchService{matchRepo: matchRepo}
}

// playerをマッチング待ち状態として保存する
func (srv *MatchService) Wait(player *model.Player) error {
	if err := srv.matchRepo.SetID(player.GetID()); err != nil {
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
