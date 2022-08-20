package model

// MatchingPlayerの状態を申請者、受諾者、状態で表す。
type MatchingState struct {
	Applicant Player
	Approver  Player
	Status    MatchingStatus
}

type MatchingStatus uint8

const (
	WAITING     MatchingStatus = iota // マッチング待ち状態
	NEGOTIATING                       // マッチング状態(対戦申請を受信または送信中)
)
