package model

type Message interface {
	GetSource() Player
	GetDest() Player
}

type MatchingMessage struct {
	source *MatchingPlayer
	dest   *MatchingPlayer
	data   *MatchingMessageData
}

func (mm *MatchingMessage) GetSource() Player {
	return mm.source
}

func (mm *MatchingMessage) GetDest() Player {
	return mm.dest
}

// TODO: マッチング時のメッセージのフォーマットの設計
type MatchingMessageData struct {
}
