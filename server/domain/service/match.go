package service

import (
	"fmt"
	"sync"
	"github.com/taise-hub/shellgame-cli/server/domain/model"
)

var (
	mu sync.Mutex
)

type MatchingService struct {
}

func NewMatchingService() *MatchingService {
	return &MatchingService{}
}

func (svc *MatchingService) changeNegotiatingState(room *model.MatchingRoom, src, dst *model.MatchingPlayer) error {
	if room.Players[src.GetID()] == nil {
		return fmt.Errorf("source player is not in the room.")
	} else if room.Players[dst.GetID()] == nil {
		return fmt.Errorf("destination player is not in the room.")
	}

	mu.Lock()
	defer mu.Unlock()

	if src.GetStatus() != model.WAITING {
		return fmt.Errorf("souce player is not WAITING")
	} else if dst.GetStatus() != model.WAITING {
		return fmt.Errorf("destination player is not WAITING")
	}

	room.Players[src.GetID()].SetStatus(model.NEGOTIATING)
	room.Players[dst.GetID()].SetStatus(model.NEGOTIATING)
	return nil
}