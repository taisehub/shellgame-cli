package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// helpModelはヘルプ画面の実装
// 現在の構想ではtea.Modelは不要なため実装してない
type helpModel struct {
	focus bool
}

func NewHelpModel() helpModel {
	return helpModel{
		focus: false,
	}
}

func (hm helpModel) Update(msg tea.Msg, tm topModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			tm.screen = ""
			return tm, screenChange("help")
		case "ctrl+c":
			return tm, tea.Quit
		}
	}
	return tm, nil
}

func (m helpModel) View() string {
	return "\n\n  here is help screen\n\n  implement me"
}
