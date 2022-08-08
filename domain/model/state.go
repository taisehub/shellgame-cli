package model

// MatchingPlayerの状態を表す。
// 実装によっては変更になる可能性あり。
type MatchingState uint8

const (
	WAITING     MatchingState = iota // マッチング待ち状態
	NEGOTIATING                      // マッチング状態(対戦申請を受信または送信中)
)
