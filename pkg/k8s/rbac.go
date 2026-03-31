// Package k8s provides Kubernetes analysis tools.
package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/g-holali-david/devkit/internal/output"
)

func CheckRBAC(kubeconfig string) error {
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		return fmt.Errorf("kubeconfig not found at %s — pass --kubeconfig or set KUBECONFIG", kubeconfig)
	}

	output.Header("Kubernetes RBAC Audit")

	fmt.Println("  This command requires a running cluster with client-go.")
	fmt.Println("  Full implementation uses k8s.io/client-go to:")
	fmt.Println()
	output.Info("List all ClusterRoleBindings with cluster-admin")
	output.Info("Find ServiceAccounts with wildcard (*) permissions")
	output.Info("Detect default ServiceAccounts used by pods")
	output.Info("Check for overly permissive RoleBindings")
	output.Info("Flag Roles that grant secrets access")
	fmt.Println()
	fmt.Println("  To fully enable, add k8s.io/client-go to go.mod and rebuild.")
	fmt.Println()

	return nil
}
