package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg struct{ err error }
func (e errMsg) Error() string { return e.err.Error() }

// 対戦、終了、ヘルプ等の画面切り替えを通知するメッセージ
type screenChangeMsg struct{}
func screenChange() tea.Cmd {
	return func() tea.Msg {
		return screenChangeMsg{}
	}
}

type MatchingMsg struct {
	Source Profile `json:"source"`
	Dest   Profile `json:"dest"`
	Data   uint8    `json:"data"`
}