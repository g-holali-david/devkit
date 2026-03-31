package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "devkit",
	Short: "DevOps CLI toolkit",
	Long: `devkit — A CLI toolkit for DevOps engineers.

Lint Dockerfiles, scaffold Helm charts, audit Kubernetes RBAC,
estimate cluster costs, and generate CI pipelines.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(dockerCmd)
	rootCmd.AddCommand(helmCmd)
	rootCmd.AddCommand(k8sCmd)
	rootCmd.AddCommand(ciCmd)
}
