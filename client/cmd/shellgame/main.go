package main

import (
	"fmt"
	"github.com/taise-hub/shellgame-cli/client/ui"
	"os"
)

func main() {
	p := ui.GetProgram()
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
