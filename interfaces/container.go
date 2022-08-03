package interfaces

import (
	"context"
	"math/rand"
	"net"
	"strconv"
	"github.com/taise-hub/shellgame-cli/domain/repository"
)

type ContainerHandler interface {
	Create(context.Context, string) (string, error)
	Exec(context.Context, string, []string) (net.Conn, error)
	Start(context.Context, string) error
	Stop(context.Context, string) error
	Remove(context.Context, string) error
}

type ContainerRepository struct {
	ContainerHandler
}

func NewContainerRepository(ch ContainerHandler) repository.ConsoleRepository {
	return &ContainerRepository { ch }
}

func (rep *ContainerRepository) StartShell() (net.Conn, error) {
	ctx := context.Background()
	name := strconv.Itoa(rand.Int())
	id, err := rep.Create(ctx, name)
	if err != nil {
		return nil, err
	}
	if err = rep.Start(ctx, id); err != nil {
		return nil, err
	}
	return rep.Exec(ctx, name, []string{"/bin/sh"})
}

// 10分以上残ってるゲーム用コンテナをストップして、削除する。
func (rep *ContainerRepository) CleanUp() error {
	panic("implement me")
}
