package stats

import (
	"sort"
	"strings"
	"time"

	"github.com/Franciss-prog/git-wrapped/internal/git"
)

// LangStat holds language usage for one file type.
type LangStat struct {
	Name       string
	Extension  string
	Count      int
	Percentage float64
}

// Summary is the aggregated yearly recap shown in the UI and export.
type Summary struct {
	TotalCommits  int
	ActiveRepos   int
	TopRepo       string
	Languages     []LangStat
	LongestStreak int
	BusiestDay    string
	PeakHour      int
	MonthlyHist   [12]int
	Persona     string
	PersonaDesc string
}

var extensionToLanguage = map[string]string{
	".c":    "C",
	".cc":   "C++",
	".cpp":  "C++",
	".cs":   "C#",
	".css":  "CSS",
	".go":   "Go",
	".h":    "C/C++",
	".hpp":  "C++",
	".html": "HTML",
	".java": "Java",
	".js":   "JavaScript",
	".jsx":  "JavaScript",
	".kt":   "Kotlin",
	".md":   "Markdown",
	".php":  "PHP",
	".py":   "Python",
	".rb":   "Ruby",
	".rs":   "Rust",
	".sh":   "Shell",
	".sql":  "SQL",
	".swift": "Swift",
	".ts":   "TypeScript",
	".tsx":  "TypeScript",
	".yaml": "YAML",
	".yml":  "YAML",
}

var weekdayNames = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

// Compute aggregates repository scan results into a yearly summary.
func Compute(repos []git.Repository) Summary {
	var summary Summary
	if len(repos) == 0 {
		summary.Persona, summary.PersonaDesc = derivePersona(-1, "", 0)
		return summary
	}

	langCounts := make(map[string]int)
	dayCounts := make(map[time.Weekday]int)
	hourCounts := make(map[int]int)
	commitsPerRepo := make(map[string]int)
	datesWithCommits := make(map[string]struct{})

	for _, repo := range repos {
		if len(repo.Commits) == 0 {
			continue
		}

		summary.ActiveRepos++
		commitsPerRepo[repo.Name] += len(repo.Commits)
		summary.TotalCommits += len(repo.Commits)

		for _, commit := range repo.Commits {
			ts := commit.Timestamp
			day := ts.Weekday()
			hour := ts.Hour()
			month := int(ts.Month()) - 1

			dayCounts[day]++
			hourCounts[hour]++
			if month >= 0 && month < 12 {
				summary.MonthlyHist[month]++
			}

			dateKey := ts.Format("2006-01-02")
			datesWithCommits[dateKey] = struct{}{}

			for _, ext := range commit.FilesChanged {
				langCounts[ext]++
			}
		}
	}

	summary.TopRepo = topRepo(commitsPerRepo)
	summary.Languages = languageStats(langCounts)
	summary.LongestStreak = longestStreak(datesWithCommits)
	summary.BusiestDay = busiestDay(dayCounts)
	summary.PeakHour = peakHour(hourCounts)
	summary.Persona, summary.PersonaDesc = derivePersona(summary.PeakHour, summary.BusiestDay, summary.LongestStreak)

	return summary
}

func topRepo(commitsPerRepo map[string]int) string {
	var top string
	maxCommits := -1

	for name, count := range commitsPerRepo {
		if count > maxCommits || (count == maxCommits && (top == "" || name < top)) {
			top = name
			maxCommits = count
		}
	}

	return top
}

func languageName(ext string) string {
	name := extensionToLanguage[ext]
	if name != "" {
		return name
	}

	stem := strings.TrimPrefix(ext, ".")
	if stem == "" {
		return "Other"
	}
	return strings.ToUpper(stem[:1]) + stem[1:]
}

func languageStats(langCounts map[string]int) []LangStat {
	if len(langCounts) == 0 {
		return nil
	}

	byName := make(map[string]int)
	primaryExt := make(map[string]string)

	total := 0
	for ext, count := range langCounts {
		name := languageName(ext)
		byName[name] += count
		total += count

		current, ok := primaryExt[name]
		if !ok || count > langCounts[current] || (count == langCounts[current] && ext < current) {
			primaryExt[name] = ext
		}
	}

	stats := make([]LangStat, 0, len(byName))
	for name, count := range byName {
		stats = append(stats, LangStat{
			Name:       name,
			Extension:  primaryExt[name],
			Count:      count,
			Percentage: float64(count) / float64(total) * 100,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Count != stats[j].Count {
			return stats[i].Count > stats[j].Count
		}
		return stats[i].Name < stats[j].Name
	})

	return stats
}

func longestStreak(dates map[string]struct{}) int {
	if len(dates) == 0 {
		return 0
	}

	sorted := make([]time.Time, 0, len(dates))
	for dateKey := range dates {
		t, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}
		sorted = append(sorted, t)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Before(sorted[j])
	})

	longest := 1
	current := 1

	for i := 1; i < len(sorted); i++ {
		prev := sorted[i-1]
		curr := sorted[i]
		if curr.Sub(prev) == 24*time.Hour {
			current++
			if current > longest {
				longest = current
			}
			continue
		}
		current = 1
	}

	return longest
}

func busiestDay(dayCounts map[time.Weekday]int) string {
	if len(dayCounts) == 0 {
		return ""
	}

	busiest := time.Sunday
	maxCount := -1

	for day := time.Sunday; day <= time.Saturday; day++ {
		count := dayCounts[day]
		if count > maxCount {
			maxCount = count
			busiest = day
		}
	}

	if maxCount <= 0 {
		return ""
	}

	return weekdayNames[busiest]
}

func peakHour(hourCounts map[int]int) int {
	if len(hourCounts) == 0 {
		return -1
	}

	peak := -1
	maxCount := -1

	for hour, count := range hourCounts {
		if count > maxCount || (count == maxCount && (peak == -1 || hour < peak)) {
			peak = hour
			maxCount = count
		}
	}

	return peak
}

func derivePersona(peakHour int, busiestDay string, longestStreak int) (string, string) {
	if peakHour < 0 && busiestDay == "" && longestStreak == 0 {
		return "Fresh Start", "A blank canvas awaits your first commit."
	}

	if longestStreak >= 14 {
		return "Consistent Craftsperson", "Day after day, commit after commit."
	}

	switch busiestDay {
	case "Saturday", "Sunday":
		return "Weekend Warrior", "You save your best code for days off."
	}

	if peakHour >= 22 || peakHour <= 4 {
		return "Midnight Hacker", "While the world sleeps, you ship."
	}

	if peakHour >= 9 && peakHour <= 17 {
		return "9-to-5 Shipper", "Professional hours, professional commits."
	}

	return "Sprint-and-Rest", "Bursts of brilliance, then time to recharge."
}
