package common

type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Message interface {
}

type MatchingMessage struct {
	Source *Profile            `json:"source"`
	Dest   *Profile            `json:"dest"`
	Data   MatchingMessageData `json:"data"`
}

type MatchingMessageData uint8

const (
	OFFER = iota + 1
	CANCEL_OFFER
	ACCEPT
	DENY
	ERROR
	JOIN
	LEAVE
)
