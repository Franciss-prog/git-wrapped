package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Franciss-prog/git-wrapped/internal/stats"
)

// Run launches the animated TUI for a prepared summary.
func Run(summary stats.Summary, year int) error {
	program := tea.NewProgram(NewModel(summary, year), tea.WithAltScreen())
	_, err := program.Run()
	return err
}
