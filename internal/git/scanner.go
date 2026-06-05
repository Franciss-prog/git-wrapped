package git

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CommitRecord holds metadata for a single commit attributed to the target user.
type CommitRecord struct {
	Hash         string
	Message      string
	Author       string
	Email        string
	Timestamp    time.Time
	FilesChanged []string
}

// Repository holds scan results for one discovered Git repository.
type Repository struct {
	Path    string
	Name    string
	Commits []CommitRecord
}

// ScanRepos walks rootDir recursively, discovers Git repositories, and returns
func ScanRepos(rootDir string, year int, userEmail string) ([]Repository, error) {
	rootDir = filepath.Clean(rootDir)

	info, err := os.Stat(rootDir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, &os.PathError{Op: "ScanRepos", Path: rootDir, Err: os.ErrInvalid}
	}

	var repos []Repository

	walkErr := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) || os.IsNotExist(err) {
				return nil
			}
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		name := d.Name()

		if name == ".git" {
			repoPath := filepath.Dir(path)
			if repo, ok := openAndScan(repoPath, year, userEmail); ok {
				repos = append(repos, repo)
			}
			return filepath.SkipDir
		}

		if strings.HasPrefix(name, ".") {
			return filepath.SkipDir
		}

		if isBareGitDir(path) {
			if repo, ok := openAndScan(path, year, userEmail); ok {
				repos = append(repos, repo)
			}
			return filepath.SkipDir
		}

		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	return repos, nil
}

func openAndScan(repoPath string, year int, userEmail string) (Repository, bool) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return Repository{}, false
	}

	commits, err := collectCommits(repo, year, userEmail)
	if err != nil {
		return Repository{}, false
	}

	return Repository{
		Path:    repoPath,
		Name:    filepath.Base(repoPath),
		Commits: commits,
	}, true
}

func collectCommits(repo *git.Repository, year int, userEmail string) ([]CommitRecord, error) {
	log, err := repo.Log(&git.LogOptions{All: true})
	if err != nil {
		return nil, nil
	}
	defer log.Close()

	var commits []CommitRecord
	for {
		commit, err := log.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if commit.Author.When.Year() != year {
			continue
		}
		if !emailsMatch(commit.Author.Email, userEmail) {
			continue
		}

		commits = append(commits, CommitRecord{
			Hash:         commit.Hash.String(),
			Message:      strings.TrimSpace(commit.Message),
			Author:       commit.Author.Name,
			Email:        commit.Author.Email,
			Timestamp:    commit.Author.When,
			FilesChanged: fileExtensionsFromCommit(commit),
		})
	}

	return commits, nil
}

func fileExtensionsFromCommit(commit *object.Commit) []string {
	stats, err := commit.Stats()
	if err != nil {
		return nil
	}

	seen := make(map[string]struct{})
	var exts []string
	for _, stat := range stats {
		ext := strings.ToLower(filepath.Ext(stat.Name))
		if ext == "" {
			continue
		}
		if _, ok := seen[ext]; ok {
			continue
		}
		seen[ext] = struct{}{}
		exts = append(exts, ext)
	}
	return exts
}

func emailsMatch(commitEmail, targetEmail string) bool {
	return strings.EqualFold(strings.TrimSpace(commitEmail), strings.TrimSpace(targetEmail))
}

func isBareGitDir(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "HEAD")); err != nil {
		return false
	}
	objects, err := os.Stat(filepath.Join(path, "objects"))
	if err != nil || !objects.IsDir() {
		return false
	}
	refs, err := os.Stat(filepath.Join(path, "refs"))
	if err != nil || !refs.IsDir() {
		return false
	}
	if _, err := os.Stat(filepath.Join(path, "config")); err != nil {
		return false
	}
	return true
}
