package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taise-hub/shellgame-cli/client/ui"
	"os"
)

func main() {
	m := ui.NewModel()
	if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
