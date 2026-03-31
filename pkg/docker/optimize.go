package docker

import (
	"fmt"
	"strings"

	"github.com/g-holali-david/devkit/internal/output"
)

func Optimize(content string) {
	lines := splitLines(content)

	output.Header("Dockerfile Optimization Suggestions")

	suggestions := 0

	// Check for multi-stage build opportunity
	fromCount := 0
	for _, line := range lines {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "FROM ") {
			fromCount++
		}
	}
	if fromCount == 1 {
		output.Warn("Consider using a multi-stage build to reduce final image size")
		fmt.Println("    Example: separate build stage from runtime stage")
		suggestions++
	}

	// Check for alpine/distroless
	hasSlim := false
	hasAlpine := false
	hasDistroless := false
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "FROM ") {
			if strings.Contains(lower, "alpine") {
				hasAlpine = true
			}
			if strings.Contains(lower, "slim") {
				hasSlim = true
			}
			if strings.Contains(lower, "distroless") {
				hasDistroless = true
			}
		}
	}
	if !hasAlpine && !hasDistroless {
		if hasSlim {
			output.Warn("Consider switching to distroless for an even smaller image")
		} else {
			output.Warn("Consider using Alpine or Distroless base images")
		}
		fmt.Println("    - Alpine: smaller (~5MB), general purpose")
		fmt.Println("    - Distroless: minimal (~2MB), no shell (more secure)")
		suggestions++
	}

	// Check layer ordering — COPY before RUN is suboptimal
	lastCopy := -1
	lastRun := -1
	for i, line := range lines {
		upper := strings.ToUpper(strings.TrimSpace(line))
		if strings.HasPrefix(upper, "COPY ") {
			lastCopy = i
		}
		if strings.HasPrefix(upper, "RUN ") {
			lastRun = i
		}
	}
	_ = lastRun

	// Check for separate RUN commands that could be merged
	consecutiveRuns := 0
	for _, line := range lines {
		upper := strings.ToUpper(strings.TrimSpace(line))
		if strings.HasPrefix(upper, "RUN ") {
			consecutiveRuns++
		}
	}
	if consecutiveRuns > 3 {
		output.Warn(fmt.Sprintf("Found %d RUN instructions — consider merging to reduce layers", consecutiveRuns))
		fmt.Println("    Use && to chain commands in a single RUN instruction")
		suggestions++
	}

	// Check for .dockerignore patterns
	for _, line := range lines {
		upper := strings.ToUpper(strings.TrimSpace(line))
		if strings.HasPrefix(upper, "COPY . ") || upper == "COPY . ." {
			output.Warn("COPY . copies everything — ensure .dockerignore is comprehensive")
			fmt.Println("    Exclude: .git, node_modules, __pycache__, *.md, tests/")
			suggestions++
			break
		}
	}

	// Check for dependency caching
	hasCopyDeps := false
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(l), "COPY ") {
			if strings.Contains(l, "requirements.txt") || strings.Contains(l, "go.mod") ||
				strings.Contains(l, "package.json") || strings.Contains(l, "Gemfile") {
				hasCopyDeps = true
			}
		}
	}
	if !hasCopyDeps && lastCopy >= 0 {
		output.Warn("Copy dependency files separately before COPY . for better caching")
		fmt.Println("    Example: COPY go.mod go.sum ./ && RUN go mod download")
		suggestions++
	}

	if suggestions == 0 {
		output.Pass("Dockerfile looks well-optimized!")
	} else {
		fmt.Printf("\n  %s suggestion(s) found\n", output.Yellow(fmt.Sprintf("%d", suggestions)))
	}
	fmt.Println()
}
