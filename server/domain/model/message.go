package model

type Message interface {
}

type MatchingMessage struct {
	Source *Profile            `json:"source"`
	Dest   *Profile            `json:"dest"`
	Data   MatchingMessageData `json:"data"`
}

type MatchingMessageData uint8

const (
	OFFER = iota
	CANCEL_OFFER
	ACCEPT
	DENY
	ERROR
	JOIN
	LEAVE
)
