package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taise-hub/shellgame-cli/common"
	"io"
)

var (
	screens = []list.Item{
		screen("対戦"),
		screen("終了"),
		screen("ヘルプ"),
	}
)

// TOP画面表示時の選択肢を扱うリスト
type screen string

func (i screen) FilterValue() string { return "" }

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

type Profile common.Profile

func (p Profile) FilterValue() string { return "" }

type profileDelegate struct{}

func (d profileDelegate) Height() int                               { return 1 }
func (d profileDelegate) Spacing() int                              { return 0 }
func (d profileDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d profileDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Profile)
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
