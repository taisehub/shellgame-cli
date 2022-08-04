package model

// Playerの状態を表す。
// 実装によっては変更になる可能性あり。
type State uint8

const (
	INACTIVE   State = iota //マッチング待ちでない状態
	ACTIVE                  // マッチング待ち状態
	WAIT_APPLY              // 申請待ち状態
	PLAY                    //対戦中の状態
)
