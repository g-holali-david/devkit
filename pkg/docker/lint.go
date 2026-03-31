// Package docker provides Dockerfile analysis tools.
package docker

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/g-holali-david/devkit/internal/output"
)

type lintRule struct {
	id      string
	check   func(lines []string, dir string) (bool, string)
	weight  int
	message string
}

func Lint(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", path, err)
	}

	lines := splitLines(string(data))
	dir := filepath.Dir(path)

	output.Header("Dockerfile Lint Report — " + path)

	rules := getRules()
	score := 100
	passed := 0
	failed := 0

	for _, rule := range rules {
		ok, detail := rule.check(lines, dir)
		if ok {
			output.Pass(rule.message)
			passed++
		} else {
			msg := rule.message
			if detail != "" {
				msg += " — " + detail
			}
			output.Fail(msg)
			score -= rule.weight
			failed++
		}
	}

	if score < 0 {
		score = 0
	}

	fmt.Println()
	output.ScoreBar(score)
	fmt.Printf("\n  %s passed, %s failed\n\n",
		output.Green(fmt.Sprintf("%d", passed)),
		output.Red(fmt.Sprintf("%d", failed)),
	)

	return nil
}

func getRules() []lintRule {
	return []lintRule{
		{
			id:      "DL001",
			weight:  10,
			message: "FROM uses a specific tag (not :latest)",
			check: func(lines []string, _ string) (bool, string) {
				for _, line := range lines {
					l := strings.TrimSpace(line)
					if !strings.HasPrefix(strings.ToUpper(l), "FROM ") {
						continue
					}
					image := strings.Fields(l)[1]
					if image == "scratch" {
						continue
					}
					if !strings.Contains(image, ":") || strings.HasSuffix(image, ":latest") {
						return false, fmt.Sprintf("image %q uses latest or no tag", image)
					}
				}
				return true, ""
			},
		},
		{
			id:      "DL002",
			weight:  15,
			message: "USER instruction sets non-root user",
			check: func(lines []string, _ string) (bool, string) {
				for _, line := range lines {
					l := strings.TrimSpace(strings.ToUpper(line))
					if strings.HasPrefix(l, "USER ") {
						user := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "USER "))
						user = strings.TrimPrefix(user, "user ")
						if user != "root" && user != "0" {
							return true, ""
						}
					}
				}
				return false, "no non-root USER found"
			},
		},
		{
			id:      "DL003",
			weight:  5,
			message: "COPY preferred over ADD for local files",
			check: func(lines []string, _ string) (bool, string) {
				for _, line := range lines {
					l := strings.TrimSpace(strings.ToUpper(line))
					if strings.HasPrefix(l, "ADD ") {
						parts := strings.Fields(line)
						if len(parts) >= 2 {
							src := parts[1]
							if !strings.HasPrefix(src, "http://") && !strings.HasPrefix(src, "https://") {
								return false, "use COPY instead of ADD for local files"
							}
						}
					}
				}
				return true, ""
			},
		},
		{
			id:      "DL004",
			weight:  10,
			message: "Multi-stage build detected",
			check: func(lines []string, _ string) (bool, string) {
				fromCount := 0
				for _, line := range lines {
					if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "FROM ") {
						fromCount++
					}
				}
				return fromCount >= 2, "single-stage build"
			},
		},
		{
			id:      "DL005",
			weight:  5,
			message: ".dockerignore file exists",
			check: func(_ []string, dir string) (bool, string) {
				_, err := os.Stat(filepath.Join(dir, ".dockerignore"))
				return err == nil, ""
			},
		},
		{
			id:      "DL006",
			weight:  5,
			message: "WORKDIR is set",
			check: func(lines []string, _ string) (bool, string) {
				for _, line := range lines {
					if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "WORKDIR ") {
						return true, ""
					}
				}
				return false, "no WORKDIR instruction"
			},
		},
		{
			id:      "DL007",
			weight:  5,
			message: "EXPOSE instruction declares ports",
			check: func(lines []string, _ string) (bool, string) {
				for _, line := range lines {
					if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "EXPOSE ") {
						return true, ""
					}
				}
				return false, ""
			},
		},
		{
			id:      "DL008",
			weight:  5,
			message: "No apt-get upgrade (prefer pinned versions)",
			check: func(lines []string, _ string) (bool, string) {
				for _, line := range lines {
					lower := strings.ToLower(line)
					if strings.Contains(lower, "apt-get upgrade") || strings.Contains(lower, "apt-get dist-upgrade") {
						return false, "apt-get upgrade can cause non-reproducible builds"
					}
				}
				return true, ""
			},
		},
		{
			id:      "DL009",
			weight:  5,
			message: "apt-get lists cleaned after install",
			check: func(lines []string, _ string) (bool, string) {
				hasApt := false
				hasClean := false
				for _, line := range lines {
					lower := strings.ToLower(line)
					if strings.Contains(lower, "apt-get install") {
						hasApt = true
					}
					if strings.Contains(lower, "rm -rf /var/lib/apt/lists") {
						hasClean = true
					}
				}
				if !hasApt {
					return true, ""
				}
				return hasClean, "apt cache not cleaned"
			},
		},
		{
			id:      "DL010",
			weight:  10,
			message: "HEALTHCHECK instruction defined",
			check: func(lines []string, _ string) (bool, string) {
				for _, line := range lines {
					if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "HEALTHCHECK ") {
						return true, ""
					}
				}
				return false, "no HEALTHCHECK (recommended for production)"
			},
		},
		{
			id:      "DL011",
			weight:  5,
			message: "pip install uses --no-cache-dir",
			check: func(lines []string, _ string) (bool, string) {
				for _, line := range lines {
					lower := strings.ToLower(line)
					if strings.Contains(lower, "pip install") && !strings.Contains(lower, "--no-cache-dir") {
						return false, "add --no-cache-dir to reduce image size"
					}
				}
				return true, ""
			},
		},
		{
			id:      "DL012",
			weight:  5,
			message: "No use of curl | sh pattern",
			check: func(lines []string, _ string) (bool, string) {
				for _, line := range lines {
					lower := strings.ToLower(line)
					if (strings.Contains(lower, "curl") || strings.Contains(lower, "wget")) &&
						(strings.Contains(lower, "| sh") || strings.Contains(lower, "| bash") || strings.Contains(lower, "|sh") || strings.Contains(lower, "|bash")) {
						return false, "piping to shell is a security risk"
					}
				}
				return true, ""
			},
		},
	}
}

func splitLines(s string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
