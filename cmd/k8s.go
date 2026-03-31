package cmd

import (
	"github.com/g-holali-david/devkit/pkg/k8s"
	"github.com/spf13/cobra"
)

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetes tools",
}

var k8sCheckRBACCmd = &cobra.Command{
	Use:   "check-rbac",
	Short: "Audit RBAC permissions in the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		return k8s.CheckRBAC(kubeconfig)
	},
}

var k8sCostCmd = &cobra.Command{
	Use:   "cost-estimate",
	Short: "Estimate cost per namespace based on resource requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		namespace, _ := cmd.Flags().GetString("namespace")
		return k8s.CostEstimate(kubeconfig, namespace)
	},
}

func init() {
	k8sCheckRBACCmd.Flags().String("kubeconfig", "", "Path to kubeconfig (default: ~/.kube/config)")
	k8sCostCmd.Flags().String("kubeconfig", "", "Path to kubeconfig (default: ~/.kube/config)")
	k8sCostCmd.Flags().StringP("namespace", "n", "", "Specific namespace (default: all)")

	k8sCmd.AddCommand(k8sCheckRBACCmd)
	k8sCmd.AddCommand(k8sCostCmd)
}
