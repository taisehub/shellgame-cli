package ui

import (
	"log"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	shellgame "github.com/taise-hub/shellgame-cli/client"
	"github.com/charmbracelet/bubbles/textinput"
)

// initModelはゲーム開始時にユーザ名を登録する画面の実装
type initModel struct {
	textInput textinput.Model
}

func NewInitModel() initModel {
	ti := textinput.New()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Focus()

	return initModel{
		textInput: ti,
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
			tm := NewTopModel()
			return tm, screenChange("init")
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
