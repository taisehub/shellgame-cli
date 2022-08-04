package service

import (
	"github.com/taise-hub/shellgame-cli/domain/model"
)

type MatchService struct {
}

// playerをマッチング待ち状態として保存する
func (s *MatchService) Wait(player *model.Player) error {
	panic("implement me")
}

// 対戦したい相手に申請を送る
func (s *MatchService) Request() {
	panic("implement me")
}

// 対戦申請の受け入れ
// playerのStateをPLAYにしてゲームを開始する
func (s *MatchService) Accept() {
	panic("implement me")
}
