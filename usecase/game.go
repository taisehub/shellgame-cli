package usecase

import (
	"github.com/taise-hub/shellgame-cli/domain/repository"
	"io"
	"log"
	"net"
)

type GameUsecase interface {
	Start(net.Conn) error
	Matching() error
}

type gameInteractor struct {
	consoleRepo repository.ConsoleRepository
}

func NewGameInteractor(consoleRepo repository.ConsoleRepository) GameUsecase {
	return &gameInteractor{
		consoleRepo: consoleRepo,
	}
}

// ゲーム開始時に利用します。
// クラアインとから受け取ったWebsocketをdocker.sockに繋げます。
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

// ゲームのマッチングをするためのメソッド
// 詳細まで考えれていないので、最終的にはもう少し細かく分かれる
func (gi *gameInteractor) Matching() error {
	panic("implement me")
}
