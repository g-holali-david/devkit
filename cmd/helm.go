package cmd

import (
	"github.com/g-holali-david/devkit/pkg/helm"
	"github.com/spf13/cobra"
)

var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Helm-related tools",
}

var helmScaffoldCmd = &cobra.Command{
	Use:   "scaffold [app-name]",
	Short: "Generate an opinionated Helm chart",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir, _ := cmd.Flags().GetString("output")
		return helm.Scaffold(args[0], outputDir)
	},
}

func init() {
	helmScaffoldCmd.Flags().StringP("output", "o", ".", "Output directory")
	helmCmd.AddCommand(helmScaffoldCmd)
}
