package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taise-hub/shellgame-cli/common"
)

// waitModelは通信待ち画面の実装
type waitModel struct {

}

func NewWaitModel() waitModel {
	return waitModel{ }
}

func (wm waitModel) Update(msg tea.Msg, mm matchModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case screenChangeMsg:
		case MatchingMsg:
			switch msg.Data {
				case common.ACCEPT:
					panic("recieve accept")
				case common.DENY:
					panic("recieve deny")
			}
		case tea.KeyMsg:
		switch msg.String() {
		case "q":
			dest, _ := mm.list.SelectedItem().(Profile)
			if dest.ID == "" {
				return mm, screenChange("wait")
			}
			mm.sendMatchingMessage(dest, common.CANCEL_OFFER)
			mm.screen = ""
			return mm, screenChange("wait")
		}
	}
	return mm, nil
}

func (wm waitModel) View() string {
	return "\n\n  通信中...\n\n"
}