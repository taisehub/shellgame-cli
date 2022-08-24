package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type helpModel struct {
	focus bool
}

func NewHelpModel() helpModel {
	return helpModel{
		focus: false,
	}
}

func (h helpModel) Update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.screen = ""
			return m, screenChange()
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m helpModel) View() string {
	return "\n\n  here is help screen\n\n  implement me"
}
