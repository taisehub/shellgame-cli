package usecase

import (
	"context"
	"time"
	"github.com/taise-hub/shellgame-cli/domain/model"
	"github.com/taise-hub/shellgame-cli/domain/repository"
	"github.com/taise-hub/shellgame-cli/domain/service"
	"io"
	"log"
	"net"
)

type GameUsecase interface {
	Start(net.Conn) error
	GetMatchingPlayers() []*model.MatchingPlayer
	WaitMatch(*model.MatchingPlayer)
}

type gameInteractor struct {
	consoleRepo  repository.ConsoleRepository
	matchService *service.MatchService
}

func NewGameInteractor(consoleRepo repository.ConsoleRepository, matchService *service.MatchService) GameUsecase {
	return &gameInteractor{
		consoleRepo:  consoleRepo,
		matchService: matchService,
	}
}

// ゲーム開始時に利用する。
// クラアインとから受け取ったコネクションをコンソールの入出力先である別のコネクションに接続する。
func (gi *gameInteractor) Start(nconn net.Conn) (err error) {
	cconn, err := gi.consoleRepo.StartShell()
	if err != nil {
		log.Printf("Error in StartShell(): %v\n", err)
		return err
	}
	defer cconn.Close()
	go func() { _, _ = io.Copy(nconn, cconn) }()
	io.Copy(cconn, nconn)
	return
}

func (gi *gameInteractor) GetMatchingPlayers() []*model.MatchingPlayer {
	mroom := model.GetMatchingRoom()
	return mroom.GetMatchingPlayers()
}

// playerをマッチング待ち状態にする。
func (gi *gameInteractor) WaitMatch(player *model.MatchingPlayer) {
	mroom := model.GetMatchingRoom()
	mroom.GetRegisterChan() <- player
	go player.ReadPump()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		player.WritePump(ctx)
		cancel()
	}()
}
