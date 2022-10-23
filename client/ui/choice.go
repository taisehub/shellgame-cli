package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taise-hub/shellgame-cli/common"
)

// 対戦要求を受けたり断ったりするときに利用するモデル
// FIXME: 名前がなんか違う。
type choiceModel struct {
	dest Profile
}

func NewChoiceModel() choiceModel {
	return choiceModel{ }
}

func (cm choiceModel) Update(msg tea.Msg, mm matchModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case screenChangeMsg:
		case MatchingMsg:
			switch msg.Data {
				case common.CANCEL_OFFER:
					// TODO: キャンセルされたことを通達する画面を挟みたい。
					mm.screen = ""
					return mm, screenChange("choice")
			}
		case tea.KeyMsg:
		switch msg.String() {
		case "y":
			mm.sendMatchingMessage(cm.dest, common.ACCEPT)
			mm.screen = "waiting"
			return mm, screenChange("choice")
		case "n":
			mm.sendMatchingMessage(cm.dest, common.DENY)
			mm.screen = "waiting"
			return mm, screenChange("choice")
		}
	}
	return mm, nil
}

func (c choiceModel) View() string {
	return "\n\n  here is Choice screen\n\n  implement me"
}