package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	tint "github.com/lrstanley/bubbletint"
)

func main() {
	tint.NewDefaultRegistry()
	tint.SetTint(tint.TintTokyoNightStorm)

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("init error: %v", err)
		os.Exit(1)
	}
}
