package ui

import (
	"os/exec"
	tea "github.com/charmbracelet/bubbletea"
	shellgame "github.com/taise-hub/shellgame-cli/client"
	"log"
)

type shellFinishedMsg struct{ err error }

func ExecShell() tea.Cmd {
	c := exec.Command("clear") // best effort
	return tea.Batch(tea.ExecProcess(c, func(err error) tea.Msg {
		return shellFinishedMsg{err}
	}),tea.Exec(&shellgame.Terminal{}, func(err error) tea.Msg {
		return shellFinishedMsg{err}
	}))
}

type matchModel struct {
	err error
}

func NewMatchModel() matchModel {
	return matchModel{}
}

func (match matchModel) Update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
		// コネクション貼ってリスト取ってくるか?
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
			match.err = msg.err
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m matchModel) View() string {
	if m.err != nil {
		log.Fatalf("Error:" + m.err.Error() + "\n")
		return "ERROR\n"
	}
	return "aaa"
}
