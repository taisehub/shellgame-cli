package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"sync"
)

var (
	lock = &sync.Mutex{}
	p    *tea.Program
)

func GetProgram() *tea.Program {
	if p == nil {
		lock.Lock()
		defer lock.Unlock()
		if p == nil { // 確実にlock.Lock()が行われることを保証する
			m := NewInitModel()
			p = tea.NewProgram(m, tea.WithAltScreen())
		}
	}
	return p
}
