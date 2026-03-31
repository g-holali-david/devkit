package docker

import (
	"testing"
)

func TestSplitLines(t *testing.T) {
	input := "FROM golang:1.22\nRUN echo hello\nUSER appuser"
	lines := splitLines(input)
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestLintRules_FROMTag(t *testing.T) {
	rules := getRules()
	fromRule := rules[0] // DL001

	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{"tagged image", []string{"FROM golang:1.22"}, true},
		{"latest tag", []string{"FROM golang:latest"}, false},
		{"no tag", []string{"FROM golang"}, false},
		{"scratch", []string{"FROM scratch"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := fromRule.check(tt.lines, "")
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLintRules_NonRootUser(t *testing.T) {
	rules := getRules()
	userRule := rules[1] // DL002

	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{"has non-root user", []string{"USER appuser"}, true},
		{"root user", []string{"USER root"}, false},
		{"no user", []string{"FROM golang:1.22"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := userRule.check(tt.lines, "")
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLintRules_MultiStage(t *testing.T) {
	rules := getRules()
	msRule := rules[3] // DL004

	single := []string{"FROM golang:1.22"}
	multi := []string{"FROM golang:1.22 AS builder", "FROM alpine:3.19"}

	if ok, _ := msRule.check(single, ""); ok {
		t.Error("single stage should fail")
	}
	if ok, _ := msRule.check(multi, ""); !ok {
		t.Error("multi stage should pass")
	}
}
