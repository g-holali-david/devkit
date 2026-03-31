package cmd

import (
	"github.com/g-holali-david/devkit/pkg/ci"
	"github.com/spf13/cobra"
)

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "CI/CD pipeline tools",
}

var ciGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a CI pipeline configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		provider, _ := cmd.Flags().GetString("provider")
		lang, _ := cmd.Flags().GetString("language")
		outputDir, _ := cmd.Flags().GetString("output")
		return ci.Generate(provider, lang, outputDir)
	},
}

func init() {
	ciGenerateCmd.Flags().StringP("provider", "p", "github", "CI provider (github, gitlab)")
	ciGenerateCmd.Flags().StringP("language", "l", "go", "Project language (go, python, node)")
	ciGenerateCmd.Flags().StringP("output", "o", ".", "Output directory")
	ciCmd.AddCommand(ciGenerateCmd)
}
