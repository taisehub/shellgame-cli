package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taise-hub/shellgame-cli/client/ui"
	"os"
)

func main() {
	m := ui.RootModel{}
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
