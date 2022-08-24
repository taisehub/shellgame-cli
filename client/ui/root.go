package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"io"
)

var (
	width             = 0
	defaultWidth      = 140
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("40"))

	title   = "シェルゲー"
	screens = []list.Item{
		screen("対戦"),
		screen("終了"),
		screen("ヘルプ"),
	}
)

// ターミナルの幅を設定
func init() {
	var err error
	width, _, err = term.GetSize(0)
	if err != nil {
		width = defaultWidth
	}
	titleStyle.Width(width).Align(lipgloss.Left)
	itemStyle.Width(width).MarginLeft(4).Align(lipgloss.Left)
	selectedItemStyle.Width(width).MarginLeft(4).Align(lipgloss.Left)
}

type screen string

func (i screen) FilterValue() string { return "" }

type screenChangeMsg struct{}

func screenChange() tea.Cmd {
	return func() tea.Msg {
		return screenChangeMsg{}
	}
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(screen)
	if !ok {
		return
	}

	str := fmt.Sprintf(" %s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render(">  " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type model struct {
	screen string
	list   list.Model
	match  matchModel
	help   helpModel
}

func NewModel() model {
	l := list.New(screens, itemDelegate{}, width, 10)
	l.Title = title
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	m := NewMatchModel()
	h := NewHelpModel()
	return model{screen: "", list: l, match: m, help: h}
}

func (m model) Init() tea.Cmd {
	return tea.Batch()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case "対戦":
		return m.match.Update(msg, m)
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
		case "q":
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(screen)
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
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.screen {
	case "対戦":
		return m.match.View()
	case "ヘルプ":
		return m.help.View()
	default:
		return m.list.View()
	}
}
