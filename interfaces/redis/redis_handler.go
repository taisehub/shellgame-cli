package interfaces

import (
	"context"
)

type RedisHandler interface {
	ListGet(context.Context, string) ([]string, error)
	ListSet(context.Context, string, []string) error
}
