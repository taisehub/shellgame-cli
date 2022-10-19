package model

import (
	"github.com/taise-hub/shellgame-cli/common"
)

type Conn interface {
	Close() error
	Write(common.Message) error
	Read(common.Message) error
}
