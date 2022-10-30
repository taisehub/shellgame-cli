package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taise-hub/shellgame-cli/common"
)

// 対戦要求を受けたり断ったりするときに利用するモデル
type matchReceivedModel struct {
	from Profile
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
			mm.sendMatchingMessage(rm.from, common.ACCEPT)
			mm.screen = "waits"
			return mm, screenChange("received")
		case "n":
			mm.sendMatchingMessage(rm.from, common.DENY)
			mm.screen = ""
			return mm, screenChange("received")
		}
	}
	return mm, nil
}

func (rm matchReceivedModel) View() string {
	return fmt.Sprintf("\n\n %sから対戦要求を受け取りました。\n\n  対戦する → y \n  断る → n\n", rm.from.Name)
}
