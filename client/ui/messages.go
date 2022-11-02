package ui

import (
	"github.com/taise-hub/shellgame-cli/common"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// 対戦、終了、ヘルプ等の画面切り替えを通知するメッセージ
type screenChangeMsg string

func screenChange(from screenChangeMsg) tea.Cmd {
	return func() tea.Msg {
		return from
	}
}

type MatchingMsg common.MatchingMessage 