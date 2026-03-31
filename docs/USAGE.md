# Guide d'utilisation — devkit CLI

## Installation

### Depuis les releases GitHub

Téléchargez le binaire correspondant à votre OS depuis la page [Releases](https://github.com/g-holali-david/devkit/releases).

```bash
# Linux / macOS
curl -sSL https://github.com/g-holali-david/devkit/releases/latest/download/devkit_linux_amd64.tar.gz | tar xz
sudo mv devkit /usr/local/bin/

# Vérifier l'installation
devkit --version
```

### Depuis les sources

```bash
go install github.com/g-holali-david/devkit@latest
```

### Build local

```bash
git clone https://github.com/g-holali-david/devkit.git
cd devkit
go build -o devkit .
./devkit --help
```

## Commandes

### Docker Lint

Analyse la qualité d'un Dockerfile avec 12 règles de bonnes pratiques.

```bash
devkit docker lint Dockerfile
devkit docker lint path/to/Dockerfile
```

Sortie exemple :
```
Dockerfile Lint Report — Dockerfile
─────────────────────────────────────

  ✓ FROM uses a specific tag (not :latest)
  ✓ USER instruction sets non-root user
  ✓ COPY preferred over ADD for local files
  ✗ Multi-stage build detected — single-stage build
  ✓ .dockerignore file exists
  ✓ WORKDIR is set
  ✓ EXPOSE instruction declares ports
  ✓ No apt-get upgrade
  ✓ apt-get lists cleaned after install
  ✗ HEALTHCHECK instruction defined — no HEALTHCHECK (recommended for production)
  ✓ pip install uses --no-cache-dir
  ✓ No use of curl | sh pattern

  ████████████████░░░░ 80/100

  10 passed, 2 failed
```

### Docker Optimize

Suggère des optimisations pour réduire la taille et améliorer la sécurité.

```bash
devkit docker optimize Dockerfile
```

Sortie exemple :
```
Dockerfile Optimization Suggestions
─────────────────────────────────────

  ⚠ Consider using a multi-stage build to reduce final image size
    Example: separate build stage from runtime stage
  ⚠ Consider switching to distroless for an even smaller image
    - Alpine: smaller (~5MB), general purpose
    - Distroless: minimal (~2MB), no shell (more secure)

  2 suggestion(s) found
```

### Helm Scaffold

Génère un Helm chart complet et opinionated.

```bash
# Génère dans le dossier courant
devkit helm scaffold my-app

# Génère dans un dossier spécifique
devkit helm scaffold my-app --output ./charts
```

Fichiers générés :
- `Chart.yaml` — métadonnées du chart
- `values.yaml` — 2 replicas, resource limits, HPA, probes, PDB
- `templates/deployment.yaml` — runAsNonRoot, resource limits, health probes
- `templates/service.yaml` — ClusterIP
- `templates/ingress.yaml` — conditionnel, nginx className
- `templates/hpa.yaml` — autoscaling CPU
- `templates/serviceaccount.yaml`
- `templates/_helpers.tpl` — labels standards K8s
- `templates/NOTES.txt` — instructions post-install

### Kubernetes RBAC Audit

```bash
# Utilise le kubeconfig par défaut (~/.kube/config)
devkit k8s check-rbac

# Spécifier un kubeconfig
devkit k8s check-rbac --kubeconfig /path/to/kubeconfig
```

### Kubernetes Cost Estimate

```bash
# Tous les namespaces
devkit k8s cost-estimate

# Un namespace spécifique
devkit k8s cost-estimate --namespace production

# Kubeconfig custom
devkit k8s cost-estimate --kubeconfig /path/to/kubeconfig -n staging
```

### CI Pipeline Generate

```bash
# GitHub Actions pour un projet Go
devkit ci generate --provider github --language go --output .

# GitLab CI pour un projet Python
devkit ci generate --provider gitlab --language python --output .

# GitHub Actions pour un projet Node.js
devkit ci generate -p github -l node -o .
```

**Providers supportés** : `github`, `gitlab`
**Langages supportés** : `go`, `python`, `node`

## Développement

### Lancer les tests

```bash
go test -v ./...
```

### Linting

```bash
golangci-lint run
```

### Ajouter une nouvelle commande

1. Créer le fichier de logique dans `pkg/<domaine>/<action>.go`
2. Créer le fichier de commande dans `cmd/<domaine>.go`
3. Ajouter la sous-commande dans le `init()` du fichier cmd
4. Écrire les tests dans `pkg/<domaine>/<action>_test.go`

### Release

```bash
# Créer un tag
git tag v0.1.0
git push origin v0.1.0

# GoReleaser génère automatiquement les binaires via GitHub Actions
```
