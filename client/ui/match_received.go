package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taise-hub/shellgame-cli/common"
)

// 対戦要求を受けたり断ったりするときに利用するモデル
type matchReceivedModel struct {
	dest Profile
}

func NewMatchRequestModel() matchReceivedModel {
	return matchReceivedModel{}
}

func (rm matchReceivedModel) Update(msg tea.Msg, mm matchModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
	case MatchingMsg:
		switch msg.Data {
		case common.CANCEL_OFFER:
			// TODO: キャンセルされたことを通達する画面を挟みたい。
			mm.screen = ""
			return mm, screenChange("received")
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			mm.sendMatchingMessage(rm.dest, common.ACCEPT)
			mm.screen = "waits"
			return mm, screenChange("received")
		case "n":
			mm.sendMatchingMessage(rm.dest, common.DENY)
			mm.screen = "waits"
			return mm, screenChange("received")
		}
	}
	return mm, nil
}

func (c matchReceivedModel) View() string {
	return "\n\n  here is Choice screen\n\n  implement me"
}
