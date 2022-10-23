package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type topModel struct {
	screen  screen
	screens list.Model
	match   matchModel
	help    helpModel
}

func NewModel() topModel {
	var m topModel

	screens := []list.Item{
		screen("対戦"),
		screen("終了"),
		screen("ヘルプ"),
	}
	s := list.New(screens, screenDelegate{}, width, 14)
	s.Title = "シェルゲー"
	s.Styles.Title = titleStyle
	s.SetShowStatusBar(false)
	s.SetFilteringEnabled(false)
	s.SetShowHelp(false)

	mm, err := NewMatchModel()
	if err != nil {
		log.Fatalf("%v\n", err.Error())
	}
	h := NewHelpModel()

	m.screen = ""
	m.screens = s
	m.match = mm
	m.help = h

	m.match.parent = &m // 子モデルであるMatchModelの親ポインタにこのモデルのアドレスを設定する
	return m
}

func (m topModel) Init() tea.Cmd {
	return tea.Batch(m.match.Init())
}

func (tm topModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch tm.screen {
	case "対戦":
		return tm.match.Update(msg)
	case "ヘルプ":
		return tm.help.Update(msg, tm)
	default:
		return updateTop(msg, tm)
	}
}

func (tm topModel) View() string {
	switch tm.screen {
	case "対戦":
		return tm.match.View()
	case "ヘルプ":
		return tm.help.View()
	default:
		return "\n" + tm.screens.View()
	}
}

func updateTop(msg tea.Msg, tm topModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "q": // TOP画面に戻る
			return tm, nil
		case "ctrl+c": // プログラムの終了
			return tm, tea.Quit
		case "enter": // 画面遷移の実行
			i, ok := tm.screens.SelectedItem().(screen)
			if !ok {
				return tm, nil
			}
			tm.screen = screen(i)
			if tm.screen == "終了" {
				return tm, tea.Quit
			}
			return tm, screenChange("top")
		}
	}
	var cmd tea.Cmd
	tm.screens, cmd = tm.screens.Update(msg)
	return tm, cmd
}
