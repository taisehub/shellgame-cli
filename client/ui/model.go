package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"io"
)

var (
	title   = "シェルゲー"
	screens = []list.Item{
		screen("対戦"),
		screen("終了"),
		screen("ヘルプ"),
	}
)

type screen string

func (i screen) FilterValue() string { return "" }

// 対戦、終了、ヘルプ等の画面切り替えを通知するメッセージ
type screenChangeMsg struct{}

// 対戦、終了、ヘルプ等の画面切り替えメッセージを通知する関数
func screenChange() tea.Cmd {
	return func() tea.Msg {
		return screenChangeMsg{}
	}
}

type screenDelegate struct{}

func (d screenDelegate) Height() int                               { return 1 }
func (d screenDelegate) Spacing() int                              { return 0 }
func (d screenDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d screenDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(screen)
	if !ok {
		return
	}

	str := fmt.Sprintf("* %s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render(">  " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

//モデルはシェルゲーのUIを表現します。
type model struct {
	screen  string
	screens list.Model
	match   matchModel
	help    helpModel
}

func NewModel() model {
	var m model

	s := list.New(screens, screenDelegate{}, width, 14)
	s.Title = title
	s.Styles.Title = titleStyle
	s.SetShowStatusBar(false)
	s.SetFilteringEnabled(false)
	s.SetShowHelp(false)

	mm := NewMatchModel()
	h := NewHelpModel()

	m.screen = ""
	m.screens = s
	m.match = mm
	m.help = h

	m.match.parent = &m // 子モデルであるMatchModelの親ポインタにこのモデルのアドレスを設定する
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch()
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
