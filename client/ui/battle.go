package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	shellgame "github.com/taise-hub/shellgame-cli/client"
	"log"
	"os/exec"
)

type shellFinishedMsg struct{ err error }

func ExecShell() tea.Cmd {
	c := exec.Command("clear") // シェルゲーコンテナに入る前に画面をクリアする
	return tea.Batch(tea.ExecProcess(c, func(err error) tea.Msg {
		return shellFinishedMsg{err}
	}), tea.Exec(&shellgame.Terminal{}, func(err error) tea.Msg {
		return shellFinishedMsg{err}
	}))
}

type battleModel struct {
	err error
}

func NewBattleModel() battleModel {
	return battleModel{}
}

func (bm battleModel) Update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
		// 問題と自分と対戦相手のスコアを取ってくる。
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, ExecShell()
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			m.screen = ""
			return m, screenChange()
		}
	case shellFinishedMsg:
		if msg.err != nil {
			bm.err = msg.err
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m battleModel) View() string {
	if m.err != nil {
		log.Fatalf("Error:" + m.err.Error() + "\n")
		return "ERROR\n"
	}
	return "aaa"
}
