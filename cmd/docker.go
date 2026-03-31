package cmd

import (
	"fmt"
	"os"

	"github.com/g-holali-david/devkit/pkg/docker"
	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker-related tools",
}

var dockerLintCmd = &cobra.Command{
	Use:   "lint [Dockerfile]",
	Short: "Lint a Dockerfile and report quality score",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return docker.Lint(args[0])
	},
}

var dockerOptimizeCmd = &cobra.Command{
	Use:   "optimize [Dockerfile]",
	Short: "Suggest optimizations for a Dockerfile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("cannot read %s: %w", args[0], err)
		}
		docker.Optimize(string(data))
		return nil
	},
}

func init() {
	dockerCmd.AddCommand(dockerLintCmd)
	dockerCmd.AddCommand(dockerOptimizeCmd)
}
