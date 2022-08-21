package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	shellgame "github.com/taise-hub/shellgame-cli/client"
)

type shellFinishedMsg struct{ err error }

func ExecShell() tea.Cmd {
	return tea.Exec(&shellgame.Terminal{}, func(err error) tea.Msg {
		return shellFinishedMsg{err}
	})
}

type RootModel struct {
	err error
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, ExecShell()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case shellFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m RootModel) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}
	return "Press 'enter' to open your SHELL.\nPress 'q' to quit.\n"
}
