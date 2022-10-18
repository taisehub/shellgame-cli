package ui

import (
	"log"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	screen  string
	screens list.Model
	match   matchModel
	help    helpModel
}

func NewModel() model {
	var m model

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

func (m model) Init() tea.Cmd {
	return tea.Batch(m.match.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case "対戦":
		return m.match.Update(msg)
	case "ヘルプ":
		return m.help.Update(msg, m)
	default:
		return updateTop(msg, m)
	}
}

func (m model) View() string {
	switch m.screen {
	case "対戦":
		return m.match.View()
	case "ヘルプ":
		return m.help.View()
	default:
		return "\n" + m.screens.View()
	}
}

func updateTop(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "q": // TOP画面に戻る
			return m, nil
		case "ctrl+c": // プログラムの終了
			return m, tea.Quit
		case "enter": // 画面遷移の実行
			i, ok := m.screens.SelectedItem().(screen)
			if !ok {
				return m, nil
			}
			m.screen = string(i)
			if m.screen == "終了" {
				return m, tea.Quit
			}
			return m, screenChange()
		}
	}
	var cmd tea.Cmd
	m.screens, cmd = m.screens.Update(msg)
	return m, cmd
}