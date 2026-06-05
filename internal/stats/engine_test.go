package stats

import (
	"testing"
	"time"

	"github.com/Franciss-prog/git-wrapped/internal/git"
)

func TestComputeNoCommits(t *testing.T) {
	summary := Compute([]git.Repository{
		{Name: "empty-repo", Commits: nil},
		{Name: "also-empty", Commits: []git.CommitRecord{}},
	})

	if summary.TotalCommits != 0 {
		t.Errorf("TotalCommits = %d, want 0", summary.TotalCommits)
	}
	if summary.ActiveRepos != 0 {
		t.Errorf("ActiveRepos = %d, want 0", summary.ActiveRepos)
	}
	if summary.TopRepo != "" {
		t.Errorf("TopRepo = %q, want empty", summary.TopRepo)
	}
	if summary.Languages != nil {
		t.Errorf("Languages = %v, want nil", summary.Languages)
	}
	if summary.LongestStreak != 0 {
		t.Errorf("LongestStreak = %d, want 0", summary.LongestStreak)
	}
	if summary.BusiestDay != "" {
		t.Errorf("BusiestDay = %q, want empty", summary.BusiestDay)
	}
	if summary.PeakHour != -1 {
		t.Errorf("PeakHour = %d, want -1", summary.PeakHour)
	}
	if summary.MonthlyHist != [12]int{} {
		t.Errorf("MonthlyHist = %v, want zeros", summary.MonthlyHist)
	}
	if summary.Persona != "Fresh Start" {
		t.Errorf("Persona = %q, want Fresh Start", summary.Persona)
	}
}

func TestComputeSingleDayAllSameHour(t *testing.T) {
	when := time.Date(2024, 6, 15, 14, 0, 0, 0, time.UTC) // Saturday

	summary := Compute([]git.Repository{{
		Name: "solo",
		Commits: []git.CommitRecord{
			{Timestamp: when, FilesChanged: []string{".go"}},
			{Timestamp: when.Add(30 * time.Minute), FilesChanged: []string{".go", ".md"}},
			{Timestamp: when.Add(2 * time.Hour), FilesChanged: []string{".py"}},
		},
	}})

	if summary.TotalCommits != 3 {
		t.Errorf("TotalCommits = %d, want 3", summary.TotalCommits)
	}
	if summary.ActiveRepos != 1 {
		t.Errorf("ActiveRepos = %d, want 1", summary.ActiveRepos)
	}
	if summary.TopRepo != "solo" {
		t.Errorf("TopRepo = %q, want solo", summary.TopRepo)
	}
	if summary.LongestStreak != 1 {
		t.Errorf("LongestStreak = %d, want 1", summary.LongestStreak)
	}
	if summary.BusiestDay != "Saturday" {
		t.Errorf("BusiestDay = %q, want Saturday", summary.BusiestDay)
	}
	if summary.PeakHour != 14 {
		t.Errorf("PeakHour = %d, want 14", summary.PeakHour)
	}
	if summary.MonthlyHist[5] != 3 {
		t.Errorf("MonthlyHist[5] = %d, want 3", summary.MonthlyHist[5])
	}
	if summary.Persona != "Weekend Warrior" {
		t.Errorf("Persona = %q, want Weekend Warrior", summary.Persona)
	}

	if len(summary.Languages) != 3 {
		t.Fatalf("Languages len = %d, want 3", len(summary.Languages))
	}
	if summary.Languages[0].Name != "Go" || summary.Languages[0].Count != 2 {
		t.Errorf("top language = %+v, want Go with count 2", summary.Languages[0])
	}
}

func TestComputeLongestStreak(t *testing.T) {
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	commits := []git.CommitRecord{
		{Timestamp: base, FilesChanged: []string{".go"}},
		{Timestamp: base.Add(24 * time.Hour), FilesChanged: []string{".go"}},
		{Timestamp: base.Add(48 * time.Hour), FilesChanged: []string{".go"}},
		{Timestamp: base.Add(96 * time.Hour), FilesChanged: []string{".go"}},
	}

	summary := Compute([]git.Repository{{Name: "streaky", Commits: commits}})

	if summary.LongestStreak != 3 {
		t.Errorf("LongestStreak = %d, want 3", summary.LongestStreak)
	}
}

func TestComputeTopRepoAndLanguages(t *testing.T) {
	when := time.Date(2024, 3, 10, 11, 0, 0, 0, time.UTC) // Sunday

	summary := Compute([]git.Repository{
		{
			Name: "alpha",
			Commits: []git.CommitRecord{
				{Timestamp: when, FilesChanged: []string{".go"}},
			},
		},
		{
			Name: "beta",
			Commits: []git.CommitRecord{
				{Timestamp: when, FilesChanged: []string{".ts"}},
				{Timestamp: when.Add(time.Hour), FilesChanged: []string{".tsx"}},
				{Timestamp: when.Add(2 * time.Hour), FilesChanged: []string{".ts"}},
			},
		},
	})

	if summary.TotalCommits != 4 {
		t.Errorf("TotalCommits = %d, want 4", summary.TotalCommits)
	}
	if summary.ActiveRepos != 2 {
		t.Errorf("ActiveRepos = %d, want 2", summary.ActiveRepos)
	}
	if summary.TopRepo != "beta" {
		t.Errorf("TopRepo = %q, want beta", summary.TopRepo)
	}
	if summary.MonthlyHist[2] != 4 {
		t.Errorf("MonthlyHist[2] = %d, want 4", summary.MonthlyHist[2])
	}

	var tsCount int
	for _, lang := range summary.Languages {
		if lang.Name == "TypeScript" {
			tsCount = lang.Count
		}
	}
	if tsCount != 3 {
		t.Errorf("TypeScript count = %d, want 3", tsCount)
	}
}

func TestComputePersonas(t *testing.T) {
	tests := []struct {
		name      string
		commits   []git.CommitRecord
		wantPersona string
	}{
		{
			name: "midnight hacker",
			commits: []git.CommitRecord{{
				Timestamp: time.Date(2024, 2, 5, 23, 0, 0, 0, time.UTC), // Monday night
				FilesChanged: []string{".go"},
			}},
			wantPersona: "Midnight Hacker",
		},
		{
			name: "9-to-5 shipper",
			commits: []git.CommitRecord{{
				Timestamp: time.Date(2024, 2, 6, 10, 0, 0, 0, time.UTC), // Tuesday morning
				FilesChanged: []string{".go"},
			}},
			wantPersona: "9-to-5 Shipper",
		},
		{
			name: "consistent craftsperson",
			commits: func() []git.CommitRecord {
				base := time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC)
				out := make([]git.CommitRecord, 14)
				for i := range out {
					out[i] = git.CommitRecord{
						Timestamp:    base.Add(time.Duration(i) * 24 * time.Hour),
						FilesChanged: []string{".go"},
					}
				}
				return out
			}(),
			wantPersona: "Consistent Craftsperson",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := Compute([]git.Repository{{Name: "persona", Commits: tt.commits}})
			if summary.Persona != tt.wantPersona {
				t.Errorf("Persona = %q, want %q", summary.Persona, tt.wantPersona)
			}
		})
	}
}

func TestComputeEmptyReposInput(t *testing.T) {
	summary := Compute(nil)
	if summary.Persona != "Fresh Start" {
		t.Errorf("Persona = %q, want Fresh Start", summary.Persona)
	}
}
