package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Franciss-prog/git-wrapped/internal/git"
)

type scanDoneMsg struct {
	repos []git.Repository
	err   error
}

type loadingModel struct {
	spinner   spinner.Model
	rootDir   string
	year      int
	userEmail string
	repos     []git.Repository
	err       error
}

func newLoadingModel(rootDir string, year int, userEmail string) loadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	return loadingModel{
		spinner:   s,
		rootDir:   rootDir,
		year:      year,
		userEmail: userEmail,
	}
}

func (m loadingModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		repos, err := git.ScanRepos(m.rootDir, m.year, m.userEmail)
		return scanDoneMsg{repos: repos, err: err}
	})
}

func (m loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.err = fmt.Errorf("scan canceled")
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		return m, nil
	case scanDoneMsg:
		m.repos = msg.repos
		m.err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m loadingModel) View() string {
	return fmt.Sprintf("%s Scanning repositories...", m.spinner.View())
}

func scanWithSpinner(rootDir string, year int, userEmail string) ([]git.Repository, error) {
	program := tea.NewProgram(newLoadingModel(rootDir, year, userEmail), tea.WithAltScreen())
	model, err := program.Run()
	if err != nil {
		return nil, err
	}

	loading, ok := model.(loadingModel)
	if !ok {
		return nil, fmt.Errorf("unexpected loading model result")
	}
	return loading.repos, loading.err
}