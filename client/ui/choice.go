package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// 対戦要求を受けたり断ったりするときに利用するモデル
// FIXME: 名前がなんか違う。
type choiceModel struct {

}

func NewChoiceModel() choiceModel {
	return choiceModel{ }
}

func (cm choiceModel) Update(msg tea.Msg, tm matchModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case screenChangeMsg:
		case tea.KeyMsg:
		switch msg.String() {
		case "q":
			tm.screen = ""
			return tm, screenChange("choice")
		}
	}
	return tm, nil
}

func (c choiceModel) View() string {
	return "\n\n  here is Choice screen\n\n  implement me"
}