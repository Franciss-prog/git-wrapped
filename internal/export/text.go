package export

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Franciss-prog/git-wrapped/internal/stats"
)

// WriteSummary writes a plain-text summary report to path.
func WriteSummary(summary stats.Summary, year int, path string) error {
	body := formatSummary(summary, year)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(body), 0o644)
}

func formatSummary(summary stats.Summary, year int) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("GitWrapped %d Summary\n", year))
	builder.WriteString(strings.Repeat("=", 32))
	builder.WriteString("\n\n")
	builder.WriteString(fmt.Sprintf("Total commits: %d\n", summary.TotalCommits))
	builder.WriteString(fmt.Sprintf("Active repos: %d\n", summary.ActiveRepos))
	if summary.TopRepo != "" {
		builder.WriteString(fmt.Sprintf("Top repo: %s\n", summary.TopRepo))
	}
	builder.WriteString(fmt.Sprintf("Longest streak: %d days\n", summary.LongestStreak))
	if summary.BusiestDay != "" {
		builder.WriteString(fmt.Sprintf("Busiest day: %s\n", summary.BusiestDay))
	}
	if summary.PeakHour >= 0 {
		builder.WriteString(fmt.Sprintf("Peak hour: %02d:00\n", summary.PeakHour))
	}
	builder.WriteString(fmt.Sprintf("Persona: %s\n", summary.Persona))
	if summary.PersonaDesc != "" {
		builder.WriteString(fmt.Sprintf("Persona note: %s\n", summary.PersonaDesc))
	}

	if len(summary.Languages) > 0 {
		builder.WriteString("\nLanguages:\n")
		langs := append([]stats.LangStat(nil), summary.Languages...)
		sort.Slice(langs, func(i, j int) bool {
			if langs[i].Count != langs[j].Count {
				return langs[i].Count > langs[j].Count
			}
			return langs[i].Name < langs[j].Name
		})
		for _, lang := range langs {
			builder.WriteString(fmt.Sprintf("- %s (%s): %d commits, %.1f%%\n", lang.Name, lang.Extension, lang.Count, lang.Percentage))
		}
	}

	builder.WriteString("\nMonthly activity:\n")
	for month, count := range summary.MonthlyHist {
		builder.WriteString(fmt.Sprintf("- %02d: %d\n", month+1, count))
	}

	return builder.String()
}