package model

type Message interface {
}

type MatchingMessage struct {
	Source uint32              `json:"source"`
	Dest   uint32              `json:"dest"`
	Data   MatchingMessageData `json:"data"`
}

type MatchingMessageData uint8

const (
	OFFER = iota
	CANCEL_OFFER
	ACCEPT
	DENY
	ERROR
)
