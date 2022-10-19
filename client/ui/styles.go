package ui

import (
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	width        = 0
	defaultWidth = 140

	baseStyle         = lipgloss.NewStyle().Margin().Padding()
	titleStyle        = baseStyle.Copy().Bold(true).Foreground(lipgloss.Color("46"))
	itemStyle         = baseStyle.Copy().Margin().Padding()
	selectedItemStyle = baseStyle.Copy().Foreground(lipgloss.Color("41"))
)

// ターミナルの幅を取得し設定する
func init() {
	var err error
	width, _, err = term.GetSize(0)
	if err != nil {
		width = defaultWidth
	}
	titleStyle.Width(width).Height(10).Align(lipgloss.Center)
	itemStyle.Width(width).Height(1).Align(lipgloss.Left).MarginLeft(width * 49 / 100)
	selectedItemStyle.Width(width).Height(1).Align(lipgloss.Left).MarginLeft(width * 49 / 100)
}
