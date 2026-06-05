package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestBarChart(t *testing.T) {
	got := BarChart("Commits", 6, 12, 6, lipgloss.Color(""))
	want := "Commits  ███░░░  6"
	if got != want {
		t.Fatalf("BarChart() = %q, want %q", got, want)
	}
}

func TestProgressBar(t *testing.T) {
	got := ProgressBar(50, 10, lipgloss.Color(""))
	want := "█████░░░░░  50.0%"
	if got != want {
		t.Fatalf("ProgressBar() = %q, want %q", got, want)
	}
}

func TestSparkline(t *testing.T) {
	values := [12]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	got := Sparkline(values, 12)
	want := "▁▂▂▃▄▄▅▅▆▇▇█"
	if got != want {
		t.Fatalf("Sparkline() = %q, want %q", got, want)
	}
}