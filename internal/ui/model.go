package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Franciss-prog/git-wrapped/internal/stats"
)

const totalSlides = 6
const maxAnimationTick = 24

type tickMsg time.Time

// Model drives the animated slide deck shown in the terminal.
type Model struct {
	Slide   int
	Summary stats.Summary
	Tick    int
	Width   int
	Height  int
	Year    int
	Done    bool
}

// NewModel creates the Bubble Tea model for the wrapped slideshow.
func NewModel(summary stats.Summary, year int) Model {
	return Model{
		Summary: summary,
		Year:    year,
		Width:   80,
		Height:  24,
	}
}

// Init starts the animation timer.
func (m Model) Init() tea.Cmd {
	return m.tickCmd()
}

// Update handles keyboard, sizing, and animation messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.Done = true
			return m, tea.Quit
		case "left":
			if m.Slide > 0 {
				m.Slide--
				m.Tick = 0
			}
			return m, m.tickCmd()
		case "right", "enter":
			if m.Slide < totalSlides-1 {
				m.Slide++
				m.Tick = 0
			}
			return m, m.tickCmd()
		}
	case tickMsg:
		if m.Tick < maxAnimationTick {
			m.Tick++
			return m, m.tickCmd()
		}
		return m, nil
	}

	return m, nil
}

// View renders the active slide.
func (m Model) View() string {
	width := clamp(m.Width, 60, 220)
	theme := newTheme()

	var body string
	switch m.Slide {
	case 0:
		body = m.renderIntro(width, theme)
	case 1:
		body = m.renderCommits(width, theme)
	case 2:
		body = m.renderLanguages(width, theme)
	case 3:
		body = m.renderTiming(width, theme)
	case 4:
		body = m.renderMonthly(width, theme)
	case 5:
		body = m.renderPersona(width, theme)
	default:
		body = m.renderIntro(width, theme)
	}

	help := theme.help.Render("←/→ or Enter to navigate · q to quit")
	frame := lipgloss.NewStyle().
		Width(width-2).
		PaddingLeft(1).
		PaddingRight(1).
		Render(body)

	return lipgloss.JoinVertical(lipgloss.Left, frame, strings.Repeat(" ", max(0, width-lipgloss.Width(help)))+help)
}

func (m Model) renderIntro(width int, theme themeSet) string {
	title := fmt.Sprintf("YOUR %d WRAPPED", m.Year)
	visible := typewriter(title, m.Tick)
	if visible == "" {
		visible = title[:1]
	}

	return panel(width, theme.title.Render(visible)+"\n\n"+theme.subtle.Render("Your coding year, wrapped in the terminal."), theme.accent)
}

func (m Model) renderCommits(width int, theme themeSet) string {
	total := m.Summary.TotalCommits
	filled := animatedFill(total, m.Tick, maxAnimationTick)
	bar := BarChart("Commits", filled, maxInt(total, 1), maxInt(18, width/3), theme.accent)
	body := strings.Join([]string{
		theme.section.Render("Commit Activity"),
		bar,
		"",
		theme.value.Render(fmt.Sprintf("%d total commits", total)),
		textOrFallback("Top repo", m.Summary.TopRepo, "No repositories discovered yet"),
	}, "\n")
	return panel(width, body, theme.accent)
}

func (m Model) renderLanguages(width int, theme themeSet) string {
	langs := m.Summary.Languages
	if len(langs) > 5 {
		langs = langs[:5]
	}

	parts := []string{theme.section.Render("Language Mix")}
	for _, lang := range langs {
		parts = append(parts, ProgressBar(lang.Percentage, maxInt(18, width-10), theme.accent))
	}
	if len(langs) == 0 {
		parts = append(parts, theme.subtle.Render("No language data yet."))
	}
	return panel(width, strings.Join(parts, "\n"), theme.accent)
}

func (m Model) renderTiming(width int, theme themeSet) string {
	peak := m.Summary.PeakHour
	clock := "--"
	if peak >= 0 {
		clock = fmt.Sprintf("%02d:00", peak)
	}
	body := strings.Join([]string{
		theme.section.Render("Streaks & Timing"),
		fmt.Sprintf("Longest streak: %s", theme.value.Render(fmt.Sprintf("%d days", m.Summary.LongestStreak))),
		fmt.Sprintf("Busiest day: %s", theme.value.Render(orFallback(m.Summary.BusiestDay, "n/a"))),
		fmt.Sprintf("Peak hour: %s", theme.value.Render(clock)),
	}, "\n")
	return panel(width, body, theme.accent)
}

func (m Model) renderMonthly(width int, theme themeSet) string {
	body := strings.Join([]string{
		theme.section.Render("Monthly Activity"),
		Sparkline(m.Summary.MonthlyHist, maxInt(12, width-6)),
	}, "\n")
	return panel(width, body, theme.accent)
}

func (m Model) renderPersona(width int, theme themeSet) string {
	name := orFallback(m.Summary.Persona, "Fresh Start")
	desc := orFallback(m.Summary.PersonaDesc, "A blank canvas awaits your first commit.")
	body := strings.Join([]string{
		theme.section.Render("Persona Reveal"),
		theme.title.Render(strings.ToUpper(name)),
		"",
		theme.subtle.Render(desc),
	}, "\n")
	return panel(width, body, theme.accent)
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func typewriter(text string, tick int) string {
	if tick <= 0 {
		return ""
	}
	if tick >= len(text) {
		return text
	}
	return text[:tick]
}

func animatedFill(final, tick, maxTick int) int {
	if final <= 0 {
		return 0
	}
	if tick >= maxTick {
		return final
	}
	filled := final * tick / maxTick
	if filled < 1 && tick > 0 {
		return 1
	}
	if filled > final {
		return final
	}
	return filled
}

func textOrFallback(label, value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return fmt.Sprintf("%s: %s", label, value)
}

func orFallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func newTheme() themeSet {
	if os.Getenv("NO_COLOR") != "" {
		plain := lipgloss.NewStyle()
		return themeSet{
			title:   plain.Bold(true),
			section: plain.Bold(true),
			value:   plain,
			subtle:  plain,
			help:    plain,
			accent:  lipgloss.Color(""),
		}
	}

	return themeSet{
		title:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
		section: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81")),
		value:   lipgloss.NewStyle().Foreground(lipgloss.Color("229")),
		subtle:  lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		help:    lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		accent:  lipgloss.Color("212"),
	}
}

type themeSet struct {
	title   lipgloss.Style
	section lipgloss.Style
	value   lipgloss.Style
	subtle  lipgloss.Style
	help    lipgloss.Style
	accent  lipgloss.Color
}

func panel(width int, body string, accent lipgloss.Color) string {
	style := lipgloss.NewStyle().
		Width(width - 2).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accent)

	return style.Render(body)
}

func clamp(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
