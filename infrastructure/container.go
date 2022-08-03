package infrastructure

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/taise-hub/shellgame-cli/interfaces"
	"log"
	"net"
	"time"
)

const (
	IMAGE_NAME = "alpine"
)

var (
	conf = &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Image:        IMAGE_NAME,
	}
	hconf = &container.HostConfig{
		AutoRemove: true,
	}
	econf = types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{},
	}
)

type containerHandler struct {
	client *client.Client
}

func NewContainerHandler() (interfaces.ContainerHandler, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	h := new(containerHandler)
	h.client = cli
	return h, nil
}

func (h *containerHandler) Create(ctx context.Context, containerName string) (string, error) {
	createdBody, err := h.client.ContainerCreate(ctx, conf, hconf, nil, nil, containerName)
	if err != nil {
		return "", err
	}
	if len(createdBody.Warnings) != 0 {
		for warn := range createdBody.Warnings {
			log.Printf("Warning in ContainerHandler.Create(): %v\n", warn)
		}
	}
	return createdBody.ID, nil
}

func (h *containerHandler) Start(ctx context.Context, containerID string) error {
	return h.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func (h *containerHandler) Exec(ctx context.Context, containerName string, cmd []string) (net.Conn, error) {
	econf.Cmd = cmd
	iresp, err := h.client.ContainerExecCreate(ctx, containerName, econf)
	if err != nil {
		return nil, err
	}
	hresp, err := h.client.ContainerExecAttach(ctx, iresp.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}
	return hresp.Conn, nil
}

func (h *containerHandler) Stop(ctx context.Context, containerID string) error {
	timeout := time.Duration(3 * time.Second)
	return h.client.ContainerStop(ctx, containerID, &timeout)
}

func (h *containerHandler) Remove(ctx context.Context, containerID string) error {
	opt := types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: true, Force: false}
	return h.client.ContainerRemove(ctx, containerID, opt)
}
