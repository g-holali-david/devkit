package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/g-holali-david/devkit/internal/output"
)

// Default hourly prices (AWS us-east-1, on-demand)
var defaultPricing = map[string]float64{
	"cpu_per_core_hour": 0.0425, // ~m5.large equivalent
	"memory_per_gb_hour": 0.0053,
}

func CostEstimate(kubeconfig, namespace string) error {
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		return fmt.Errorf("kubeconfig not found at %s — pass --kubeconfig or set KUBECONFIG", kubeconfig)
	}

	output.Header("Kubernetes Cost Estimation")

	scope := "all namespaces"
	if namespace != "" {
		scope = "namespace: " + namespace
	}
	fmt.Printf("  Scope: %s\n\n", output.Cyan(scope))

	fmt.Println("  This command requires a running cluster with client-go.")
	fmt.Println("  Full implementation reads resource requests per namespace and calculates:")
	fmt.Println()
	output.Info(fmt.Sprintf("CPU cost:    $%.4f / core / hour", defaultPricing["cpu_per_core_hour"]))
	output.Info(fmt.Sprintf("Memory cost: $%.4f / GB / hour", defaultPricing["memory_per_gb_hour"]))
	fmt.Println()
	fmt.Println("  Output format:")
	fmt.Println("  ┌─────────────┬──────┬────────┬───────────┬────────────┐")
	fmt.Println("  │ Namespace   │ CPUs │ Memory │ $/hour    │ $/month    │")
	fmt.Println("  ├─────────────┼──────┼────────┼───────────┼────────────┤")
	fmt.Println("  │ production  │ 4.0  │ 8Gi    │ $0.2124   │ $153       │")
	fmt.Println("  │ staging     │ 2.0  │ 4Gi    │ $0.1062   │ $76        │")
	fmt.Println("  └─────────────┴──────┴────────┴───────────┴────────────┘")
	fmt.Println()
	fmt.Println("  To fully enable, add k8s.io/client-go to go.mod and rebuild.")
	fmt.Println()

	return nil
}
