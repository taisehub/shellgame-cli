package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	shellgame "github.com/taise-hub/shellgame-cli/client"
	"log"
)

// initModelはゲーム開始時にユーザ名を登録する画面の実装
type initModel struct {
	textInput textinput.Model
	top       *topModel
}

func NewInitModel() initModel {
	ti := textinput.New()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Focus()
	tm := NewTopModel()

	return initModel{
		textInput: ti,
		top:       &tm,
	}
}

func (im initModel) Init() tea.Cmd {
	return nil
}

func (im initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case screenChangeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return im, tea.Quit
		case "enter":
			if err := shellgame.PostProfile(im.textInput.Value()); err != nil {
				log.Fatalf("%v\n", err.Error())
			}
			return im.top, screenChange("init")
		}
	}
	im.textInput, cmd = im.textInput.Update(msg)
	return im, cmd
}

func (im initModel) View() string {
	return fmt.Sprintf(
		"What’s your name?\n\n%s\n\n%s",
		im.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
