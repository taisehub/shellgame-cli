package usecase

import (
	"context"
	"github.com/taise-hub/shellgame-cli/common"
	"github.com/taise-hub/shellgame-cli/server/domain/model"
	"github.com/taise-hub/shellgame-cli/server/domain/repository"
	"github.com/taise-hub/shellgame-cli/server/domain/service"
	"io"
	"log"
	"net"
	"time"
)

type GameUsecase interface {
	Start(net.Conn) error
	GetMatchingProfiles() []*common.Profile
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

func (gi *gameInteractor) GetMatchingProfiles() []*common.Profile {
	mroom := model.GetMatchingRoom()
	players := mroom.GetMatchingPlayers()
	var profiles []*common.Profile
	for _, v := range players {
		profiles = append(profiles, v.Profile)
	}
	return profiles
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
