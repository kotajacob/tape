package main

import (
	"fmt"
	"os"

	"git.sr.ht/~kota/tape/tape"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	t := tape.New()
	p := tea.NewProgram(t, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error running program: %v\n",
			err,
		)
		os.Exit(1)
	}
}
