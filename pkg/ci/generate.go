// Package ci provides CI pipeline generation tools.
package ci

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/g-holali-david/devkit/internal/output"
)

func Generate(provider, lang, outputDir string) error {
	output.Header(fmt.Sprintf("Generating %s CI pipeline for %s", provider, lang))

	var content string
	var filename string

	switch provider {
	case "github":
		filename = filepath.Join(outputDir, ".github", "workflows", "ci.yml")
		content = githubWorkflow(lang)
	case "gitlab":
		filename = filepath.Join(outputDir, ".gitlab-ci.yml")
		content = gitlabCI(lang)
	default:
		return fmt.Errorf("unsupported provider: %s (supported: github, gitlab)", provider)
	}

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}

	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	output.Pass("Generated " + filename)
	fmt.Println()
	return nil
}

func githubWorkflow(lang string) string {
	switch lang {
	case "go":
		return `name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: actions/upload-artifact@v4
        with:
          name: coverage
          path: coverage.out

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: golangci-lint
        uses: golangci-lint/golangci-lint-action@v6

  build:
    runs-on: ubuntu-latest
    needs: [test, lint]
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: Build
        run: go build -o bin/ ./...
`

	case "python":
		return `name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-python@v5
        with:
          python-version: "3.12"

      - name: Install dependencies
        run: |
          pip install -r requirements.txt
          pip install pytest pytest-cov ruff

      - name: Lint
        run: ruff check .

      - name: Test
        run: pytest --cov=. --cov-report=xml

      - name: Upload coverage
        uses: actions/upload-artifact@v4
        with:
          name: coverage
          path: coverage.xml
`

	case "node":
		return `name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: "20"
          cache: "npm"

      - run: npm ci
      - run: npm run lint
      - run: npm test
      - run: npm run build
`

	default:
		return fmt.Sprintf("# CI pipeline for %s — customize as needed\n", lang)
	}
}

func gitlabCI(lang string) string {
	switch lang {
	case "go":
		return `stages:
  - test
  - build

test:
  stage: test
  image: golang:1.22
  script:
    - go test -v -race ./...

lint:
  stage: test
  image: golangci/golangci-lint:latest
  script:
    - golangci-lint run

build:
  stage: build
  image: golang:1.22
  script:
    - go build -o bin/ ./...
  artifacts:
    paths:
      - bin/
`

	case "python":
		return `stages:
  - test

test:
  stage: test
  image: python:3.12
  script:
    - pip install -r requirements.txt
    - pip install pytest ruff
    - ruff check .
    - pytest
`

	default:
		return fmt.Sprintf("# GitLab CI pipeline for %s\nstages:\n  - test\n  - build\n", lang)
	}
}
