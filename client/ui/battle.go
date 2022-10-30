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
	isShell bool
	err error
}

func NewBattleModel() battleModel {
	return battleModel{isShell: false}
}

func (bm battleModel) Init() tea.Cmd {
	return nil
}

func (bm battleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
		// 問題と自分と対戦相手のスコアを取ってくる。
		return bm, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			bm.isShell = true
			return bm, ExecShell()
		case "ctrl+c":
			return bm, tea.Quit
		}
	case shellFinishedMsg:
		bm.isShell = false
		if msg.err != nil {
			bm.err = msg.err
			return bm, tea.Quit
		}
	}
	return bm, nil
}

func (bm battleModel) View() string {
	if bm.err != nil {
		log.Fatalf("Error:" + bm.err.Error() + "\n")
		return "ERROR\n"
	}

	if bm.isShell {
		return ""
	}
	return "対戦画面\n\nIMPLEMENT ME\n\n(enterでコンテナに接続します)"
}
