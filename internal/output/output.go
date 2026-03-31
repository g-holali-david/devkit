// Package output provides colored terminal output helpers.
package output

import "fmt"

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
	bold   = "\033[1m"
	dim    = "\033[2m"
)

func Red(s string) string    { return red + s + reset }
func Green(s string) string  { return green + s + reset }
func Yellow(s string) string { return yellow + s + reset }
func Blue(s string) string   { return blue + s + reset }
func Cyan(s string) string   { return cyan + s + reset }
func Bold(s string) string   { return bold + s + reset }
func Dim(s string) string    { return dim + s + reset }

func Pass(label string)  { fmt.Printf("  %s %s\n", Green("✓"), label) }
func Warn(label string)  { fmt.Printf("  %s %s\n", Yellow("⚠"), label) }
func Fail(label string)  { fmt.Printf("  %s %s\n", Red("✗"), label) }
func Info(label string)  { fmt.Printf("  %s %s\n", Blue("ℹ"), label) }
func Score(score int)    { fmt.Printf("\n  %s %d/100\n", Bold("Score:"), score) }

func Header(title string) {
	fmt.Printf("\n%s\n%s\n\n", Bold(title), dim+"─────────────────────────────────────"+reset)
}

func ScoreBar(score int) {
	color := red
	switch {
	case score >= 80:
		color = green
	case score >= 60:
		color = yellow
	}

	filled := score / 5
	empty := 20 - filled
	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}
	fmt.Printf("  %s%s%s %d/100\n", color, bar, reset, score)
}
