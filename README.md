# GitWrapped

> _Your coding year, wrapped — right in your terminal._

GitWrapped is a terminal-based developer recap tool built with **Go** and **BubbleTea**. It scans your local Git repositories and turns your coding activity into a fun, animated yearly summary — no browser required.

---

## What It Does

Inspired by Spotify Wrapped, GitWrapped digs through your commit history and presents your development journey as a slick, interactive terminal experience. Think stats, streaks, personality summaries, and ASCII charts — all without leaving your shell.

---

## Features

- **Git Commit Analysis** — Deep-dives into your commit history across repos
- **Repository Activity Tracking** — See which projects got the most love
- **Language Usage Statistics** — Find out what you _actually_ code in
- **Coding Streaks & Productivity Insights** — Longest streaks, busiest days, peak hours
- **Developer Personality Summaries** — Are you a _Midnight Hacker_ or a _9-to-5 Shipper_?
- **Interactive Terminal UI** — Smooth animations powered by BubbleTea
- **ASCII Charts & Visualizations** — Beautiful data, no GUI needed

---

## Tech Stack

| Layer         | Technology                                              |
| ------------- | ------------------------------------------------------- |
| Language      | Go                                                      |
| TUI Framework | [BubbleTea](https://github.com/charmbracelet/bubbletea) |
| Styling       | [Lip Gloss](https://github.com/charmbracelet/lipgloss)  |
| Charts        | ASCII / ANSI rendering                                  |
| Data Source   | Local `.git` directories                                |

---

## Installation

```bash
# Clone the repo
git clone https://github.com/Franciss-prog/git-wrapped.git
cd git-wrapped

# Build
go build -o git-wrapped .

# Run
./git-wrapped
```

> **Requirements:** Go 1.21+ and at least one local Git repository to analyze.

---

## Usage

```bash
# Analyze the current year (default)
./git-wrapped

# Analyze a specific year
./git-wrapped --year 2023

# Point to a custom repos directory
./git-wrapped --dir ~/projects

# Export your summary as text
./git-wrapped --export summary.txt
```

When `--export` is set, GitWrapped writes a plain-text summary and exits without launching the terminal UI.

---

## File Structure

```
git-wrapped/
│
├── cmd/
│ └── git-wrapped/
│ └── main.go
│
├── internal/
│ ├── git/
│ ├── stats/
│ ├── tui/
│ └── export/
│
├── README.md
├── LICENSE
├── go.mod
└── go.sum
```

---

## Preview

```
+==========================================+
| YOUR 2024 WRAPPED |
+==========================================+
| Total Commits ████████████ 1,337 |
| Active Repos ██████░░░░░░ 12 |
| Top Language Go 68% |
| Longest Streak ░░░░░░░░░░░░ 21 days |
| |
| You are a: MIDNIGHT HACKER |
| Peak coding hour: 11 PM - 2 AM |
+==========================================+

```

---

````

## Roadmap

- [ ] Multi-year comparison mode
- [ ] Team/org-wide summaries
- [ ] GitHub remote repo support
- [ ] Shareable summary cards (exported as styled text)
- [ ] Plugin system for custom stat modules

---

## Contributing

Contributions are welcome! Please open an issue first to discuss what you'd like to change.

```bash
# Fork -> branch -> commit -> PR
git checkout -b feature/your-idea
````

---

## License

MIT © [7\](https://github.com/Franciss-prog)

---

<p align="center">
  Built with Go and too many late-night commits.
</p>
