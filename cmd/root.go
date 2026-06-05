package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Franciss-prog/git-wrapped/internal/export"
	"github.com/Franciss-prog/git-wrapped/internal/stats"
	"github.com/Franciss-prog/git-wrapped/internal/ui"
)

const Version = "v0.1.0"

// Run parses CLI flags and starts the application.
func Run() error {
	yearFlag := flag.Int("year", time.Now().Year(), "year to analyze (default: current year)")
	dirFlag := flag.String("dir", "~", "root directory to scan (default: home)")
	exportFlag := flag.String("export", "", "export summary to a text file")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("git-wrapped %s\n", Version)
		return nil
	}

	dir, err := resolveDir(*dirFlag)
	if err != nil {
		return err
	}

	userEmail, err := gitUserEmail()
	if err != nil {
		return err
	}

	repos, err := scanWithSpinner(dir, *yearFlag, userEmail)
	if err != nil {
		return err
	}

	summary := stats.Compute(repos)
	if strings.TrimSpace(*exportFlag) != "" {
		if err := export.WriteSummary(summary, *yearFlag, *exportFlag); err != nil {
			return err
		}
		fmt.Printf("exported summary to %s\n", *exportFlag)
		return nil
	}

	return ui.Run(summary, *yearFlag)
}

func resolveDir(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			path = home
		} else {
			path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("scan directory %q: %w", abs, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("scan directory %q is not a directory", abs)
	}
	return abs, nil
}

func gitUserEmail() (string, error) {
	cmd := exec.Command("git", "config", "--global", "user.email")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("read git user.email: %w", err)
	}
	email := strings.TrimSpace(string(output))
	if email == "" {
		return "", fmt.Errorf("git user.email is not configured")
	}
	return email, nil
}
