package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	shellgame "github.com/taise-hub/shellgame-cli/client"
	"github.com/charmbracelet/bubbles/list"
)

type shellFinishedMsg struct{ err error }

func ExecShell() tea.Cmd {
	return tea.Exec(&shellgame.Terminal{}, func(err error) tea.Msg {
		return shellFinishedMsg{err}
	})
}

type battleModel struct {
	screen screen
	screens list.Model
	isShell bool
	err error
}

func NewBattleModel() battleModel {
	screens := []list.Item{
		screen("シェル"),
		screen("回答送信"),
		screen("降参"),
	}
	s := list.New(screens, screenDelegate{}, width, 14)
	s.Title = "対戦中..."
	s.Styles.Title = titleStyle
	s.SetShowStatusBar(false)
	s.SetFilteringEnabled(false)
	s.SetShowHelp(false)

	return battleModel{isShell: false, screen: "", screens: s}
}

func (bm battleModel) Init() tea.Cmd {
	return nil
}

func (bm battleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
		// 問題と自分と対戦相手のスコアを取ってくる。
		return bm, nil
	case shellFinishedMsg:
		bm.screen = screen("")
		if msg.err != nil {
			bm.err = msg.err
			return bm, tea.Quit
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := bm.screens.SelectedItem().(screen)
			if !ok {
				return bm, nil
			}
			bm.screen = screen(i)

			switch bm.screen {
			case "シェル":
				return bm, ExecShell()
			case "降参":
				//TODO:　降参した旨を送信する
				return bm, tea.Quit
			}
		}
	}
	var cmd tea.Cmd
	bm.screens, cmd = bm.screens.Update(msg)
	return bm, cmd
}

func (bm battleModel) View() string {
	if bm.err != nil {
		panic(bm.err)
	}

	switch bm.screen {
	case "シェル":
		return ""
	default:
		return "\n" + bm.screens.View()
	}
}
