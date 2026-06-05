package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var sparkBlocks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// BarChart renders a single labeled bar as plain text.
func BarChart(label string, value, maxValue int, width int, color lipgloss.Color) string {
	filled := progressWidth(value, maxValue, width)
		bar := strings.Repeat("█", filled) + strings.Repeat("░", maxInt(0, width-filled))
	if string(color) != "" {
		bar = lipgloss.NewStyle().Foreground(color).Render(bar)
	}
	return fmt.Sprintf("%s  %s  %d", label, bar, value)
}

// Sparkline renders twelve monthly buckets using block characters.
func Sparkline(values [12]int, width int) string {
	maxValue := 0
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	if maxValue == 0 {
		return strings.Repeat(string(sparkBlocks[0]), minInt(width, 12))
	}

	var builder strings.Builder
	for _, value := range values {
		index := int(math.Round((float64(value) / float64(maxValue)) * float64(len(sparkBlocks)-1)))
		if index < 0 {
			index = 0
		}
		if index >= len(sparkBlocks) {
			index = len(sparkBlocks) - 1
		}
		builder.WriteRune(sparkBlocks[index])
	}
	return builder.String()
}

// ProgressBar renders a percentage bar for language distribution.
func ProgressBar(pct float64, width int, color lipgloss.Color) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	filled := int(math.Round((pct / 100.0) * float64(width)))
		bar := strings.Repeat("█", filled) + strings.Repeat("░", maxInt(0, width-filled))
	if string(color) != "" {
		bar = lipgloss.NewStyle().Foreground(color).Render(bar)
	}
	return fmt.Sprintf("%s %5.1f%%", bar, pct)
}

func progressWidth(value, maxValue, width int) int {
	if maxValue <= 0 || width <= 0 {
		return 0
	}
	filled := int(math.Round((float64(value) / float64(maxValue)) * float64(width)))
	if filled < 0 {
		return 0
	}
	if filled > width {
		return width
	}
	return filled
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
