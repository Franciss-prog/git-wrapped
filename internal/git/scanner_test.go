package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const (
	testUserEmail = "scanner-test@example.com"
	testUserName  = "Scanner Test"
	testYear      = 2024
)

func TestScanReposBareRepo(t *testing.T) {
	root := t.TempDir()
	barePath := filepath.Join(root, "fixture.git")
	if err := initBareRepoFixture(barePath); err != nil {
		t.Fatalf("init bare repo fixture: %v", err)
	}

	repos, err := ScanRepos(root, testYear, testUserEmail)
	if err != nil {
		t.Fatalf("ScanRepos: %v", err)
	}

	if len(repos) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(repos))
	}

	repo := repos[0]
	if repo.Name != "fixture.git" {
		t.Errorf("Name = %q, want fixture.git", repo.Name)
	}
	if repo.Path != barePath {
		t.Errorf("Path = %q, want %q", repo.Path, barePath)
	}
	if len(repo.Commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(repo.Commits))
	}

	commit := repo.Commits[0]
	if commit.Email != testUserEmail {
		t.Errorf("Email = %q, want %q", commit.Email, testUserEmail)
	}
	if commit.Author != testUserName {
		t.Errorf("Author = %q, want %q", commit.Author, testUserName)
	}
	if commit.Message != "add main.go" {
		t.Errorf("Message = %q, want add main.go", commit.Message)
	}
	if commit.Timestamp.Year() != testYear {
		t.Errorf("Timestamp year = %d, want %d", commit.Timestamp.Year(), testYear)
	}
	if len(commit.FilesChanged) != 1 || commit.FilesChanged[0] != ".go" {
		t.Errorf("FilesChanged = %v, want [.go]", commit.FilesChanged)
	}
	if commit.Hash == "" {
		t.Error("expected non-empty commit hash")
	}
}

func TestScanReposFiltersByYearAndEmail(t *testing.T) {
	root := t.TempDir()
	barePath := filepath.Join(root, "filtered.git")
	if err := initBareRepoFixtureWithExtraCommits(barePath); err != nil {
		t.Fatalf("init bare repo fixture: %v", err)
	}

	repos, err := ScanRepos(root, testYear, testUserEmail)
	if err != nil {
		t.Fatalf("ScanRepos: %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(repos))
	}
	if len(repos[0].Commits) != 1 {
		t.Fatalf("expected 1 filtered commit, got %d", len(repos[0].Commits))
	}
}

func TestScanReposEmptyRepo(t *testing.T) {
	root := t.TempDir()
	barePath := filepath.Join(root, "empty.git")
	if err := runGit(root, "init", "--bare", barePath); err != nil {
		t.Fatalf("init empty bare repo: %v", err)
	}

	repos, err := ScanRepos(root, testYear, testUserEmail)
	if err != nil {
		t.Fatalf("ScanRepos: %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(repos))
	}
	if len(repos[0].Commits) != 0 {
		t.Fatalf("expected 0 commits, got %d", len(repos[0].Commits))
	}
}

func TestScanReposSkipsHiddenDirectories(t *testing.T) {
	root := t.TempDir()
	hidden := filepath.Join(root, ".hidden-project")
	if err := os.MkdirAll(hidden, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := runGit(hidden, "init"); err != nil {
		t.Fatalf("init hidden repo: %v", err)
	}

	repos, err := ScanRepos(root, testYear, testUserEmail)
	if err != nil {
		t.Fatalf("ScanRepos: %v", err)
	}
	if len(repos) != 0 {
		t.Fatalf("expected hidden repo to be skipped, got %d repos", len(repos))
	}
}

func TestScanReposStandardRepo(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "my-project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := runGit(project, "init"); err != nil {
		t.Fatal(err)
	}
	if err := configureGitUser(project); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, "readme.md"), []byte("# hi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runGit(project, "add", "readme.md"); err != nil {
		t.Fatal(err)
	}
	if err := runGitWithEnv(project, map[string]string{
		"GIT_AUTHOR_DATE":    "2024-03-10T10:00:00",
		"GIT_COMMITTER_DATE": "2024-03-10T10:00:00",
	}, "commit", "-m", "docs"); err != nil {
		t.Fatal(err)
	}

	repos, err := ScanRepos(root, testYear, testUserEmail)
	if err != nil {
		t.Fatalf("ScanRepos: %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(repos))
	}
	if repos[0].Name != "my-project" {
		t.Errorf("Name = %q, want my-project", repos[0].Name)
	}
	if len(repos[0].Commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(repos[0].Commits))
	}
	if repos[0].Commits[0].FilesChanged[0] != ".md" {
		t.Errorf("FilesChanged = %v, want [.md]", repos[0].Commits[0].FilesChanged)
	}
}

func initBareRepoFixture(barePath string) error {
	root := filepath.Dir(barePath)
	workPath := filepath.Join(root, ".work")

	if err := runGit(root, "init", "--bare", barePath); err != nil {
		return err
	}
	if err := runGit(root, "init", workPath); err != nil {
		return err
	}
	if err := configureGitUser(workPath); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(workPath, "main.go"), []byte("package main\n"), 0o644); err != nil {
		return err
	}
	if err := runGit(workPath, "add", "main.go"); err != nil {
		return err
	}
	if err := runGitWithEnv(workPath, map[string]string{
		"GIT_AUTHOR_DATE":    "2024-06-15T12:00:00",
		"GIT_COMMITTER_DATE": "2024-06-15T12:00:00",
	}, "commit", "-m", "add main.go"); err != nil {
		return err
	}
	if err := runGit(workPath, "branch", "-M", "main"); err != nil {
		return err
	}
	if err := runGit(workPath, "remote", "add", "origin", barePath); err != nil {
		return err
	}
	return runGit(workPath, "push", "-u", "origin", "main")
}

func initBareRepoFixtureWithExtraCommits(barePath string) error {
	if err := initBareRepoFixture(barePath); err != nil {
		return err
	}

	root := filepath.Dir(barePath)
	workPath := filepath.Join(root, ".work")

	if err := os.WriteFile(filepath.Join(workPath, "other.txt"), []byte("other\n"), 0o644); err != nil {
		return err
	}
	if err := runGit(workPath, "add", "other.txt"); err != nil {
		return err
	}
	if err := runGitWithEnv(workPath, map[string]string{
		"GIT_AUTHOR_DATE":    "2023-01-01T00:00:00",
		"GIT_COMMITTER_DATE": "2023-01-01T00:00:00",
	}, "commit", "-m", "old commit"); err != nil {
		return err
	}

	if err := configureGitUserEmail(workPath, "other@example.com"); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(workPath, "third.rs"), []byte("fn main() {}\n"), 0o644); err != nil {
		return err
	}
	if err := runGit(workPath, "add", "third.rs"); err != nil {
		return err
	}
	if err := runGitWithEnv(workPath, map[string]string{
		"GIT_AUTHOR_DATE":    time.Date(testYear, 7, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		"GIT_COMMITTER_DATE": time.Date(testYear, 7, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
	}, "commit", "-m", "wrong author"); err != nil {
		return err
	}

	return runGit(workPath, "push", "origin", "main")
}

func configureGitUser(dir string) error {
	if err := runGit(dir, "config", "user.email", testUserEmail); err != nil {
		return err
	}
	return runGit(dir, "config", "user.name", testUserName)
}

func configureGitUserEmail(dir, email string) error {
	return runGit(dir, "config", "user.email", email)
}

func runGit(dir string, args ...string) error {
	return runGitWithEnv(dir, nil, args...)
}

func runGitWithEnv(dir string, extraEnv map[string]string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	for key, value := range extraEnv {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &gitCommandError{args: args, output: string(out), err: err}
	}
	return nil
}

type gitCommandError struct {
	args   []string
	output string
	err    error
}

func (e *gitCommandError) Error() string {
	return "git " + joinArgs(e.args) + ": " + e.err.Error() + "\n" + e.output
}

func joinArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}
	result := args[0]
	for _, arg := range args[1:] {
		result += " " + arg
	}
	return result
}
