package export

import (
	"github.com/Franciss-prog/git-wrapped/internal/stats"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteSummary(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "summary.txt")

	summary := stats.Summary{
		TotalCommits:  12,
		ActiveRepos:   2,
		TopRepo:       "git-wrapped",
		LongestStreak: 7,
		BusiestDay:    "Tuesday",
		PeakHour:      23,
		Persona:       "Midnight Hacker",
		PersonaDesc:   "While the world sleeps, you ship.",
		Languages: []stats.LangStat{
			{Name: "Go", Extension: ".go", Count: 8, Percentage: 66.7},
			{Name: "Markdown", Extension: ".md", Count: 4, Percentage: 33.3},
		},
		MonthlyHist: [12]int{1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	if err := WriteSummary(summary, 2024, path); err != nil {
		t.Fatalf("WriteSummary() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	wantContains := []string{
		"GitWrapped 2024 Summary",
		"Total commits: 12",
		"Active repos: 2",
		"Top repo: git-wrapped",
		"Persona: Midnight Hacker",
		"- Go (.go): 8 commits, 66.7%",
		"- Markdown (.md): 4 commits, 33.3%",
		"- 01: 1",
		"- 02: 2",
	}
	content := string(got)
	for _, fragment := range wantContains {
		if !strings.Contains(content, fragment) {
			t.Fatalf("summary output missing %q in:\n%s", fragment, content)
		}
	}
}

