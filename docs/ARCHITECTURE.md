# Architecture — devkit CLI

## Vue d'ensemble

`devkit` est un outil CLI écrit en Go qui automatise les tâches DevOps récurrentes. Il utilise le pattern **Cobra** pour la gestion des commandes avec une architecture modulaire `cmd/` + `pkg/`.

## Structure du projet

```
devkit/
├── main.go                    # Point d'entrée
├── cmd/                       # Commandes Cobra (interface CLI)
│   ├── root.go                #   Commande racine + version
│   ├── docker.go              #   devkit docker [lint|optimize]
│   ├── helm.go                #   devkit helm [scaffold]
│   ├── k8s.go                 #   devkit k8s [check-rbac|cost-estimate]
│   └── ci.go                  #   devkit ci [generate]
├── pkg/                       # Logique métier (réutilisable)
│   ├── docker/
│   │   ├── lint.go            #   12 règles de linting Dockerfile
│   │   ├── lint_test.go       #   Tests unitaires
│   │   └── optimize.go        #   Suggestions d'optimisation
│   ├── helm/
│   │   └── scaffold.go        #   Génération de Helm charts
│   ├── k8s/
│   │   ├── rbac.go            #   Audit RBAC (nécessite client-go)
│   │   └── cost.go            #   Estimation de coûts
│   └── ci/
│       └── generate.go        #   Génération de pipelines CI
├── internal/
│   └── output/
│       └── output.go          #   Helpers d'affichage terminal coloré
├── .goreleaser.yml            # Config GoReleaser (builds multi-OS)
└── .github/workflows/
    ├── ci.yml                 # Tests + lint + build
    └── release.yml            # Release automatique sur tag
```

## Principes d'architecture

### Séparation cmd / pkg

- **`cmd/`** : uniquement la définition des commandes Cobra, le parsing des flags, et l'appel aux fonctions `pkg/`
- **`pkg/`** : toute la logique métier, testable indépendamment, réutilisable comme bibliothèque Go
- **`internal/`** : helpers internes (pas exportés en dehors du module)

### Flux d'exécution

```
main.go
  └── cmd.Execute()
        └── rootCmd.Execute() (Cobra)
              ├── docker lint <file>    → pkg/docker.Lint(path)
              ├── docker optimize <file> → pkg/docker.Optimize(content)
              ├── helm scaffold <name>  → pkg/helm.Scaffold(name, outputDir)
              ├── k8s check-rbac       → pkg/k8s.CheckRBAC(kubeconfig)
              ├── k8s cost-estimate     → pkg/k8s.CostEstimate(kubeconfig, ns)
              └── ci generate           → pkg/ci.Generate(provider, lang, dir)
```

## Commandes détaillées

### `devkit docker lint <Dockerfile>`

Analyse un Dockerfile et produit un score de qualité sur 100.

**12 règles de linting :**

| ID | Poids | Règle |
|----|-------|-------|
| DL001 | 10 | FROM utilise un tag spécifique (pas :latest) |
| DL002 | 15 | Instruction USER avec utilisateur non-root |
| DL003 | 5 | COPY préféré à ADD pour fichiers locaux |
| DL004 | 10 | Build multi-stage détecté |
| DL005 | 5 | Fichier .dockerignore présent |
| DL006 | 5 | WORKDIR défini |
| DL007 | 5 | Instruction EXPOSE déclare les ports |
| DL008 | 5 | Pas de apt-get upgrade |
| DL009 | 5 | Cache apt nettoyé après install |
| DL010 | 10 | Instruction HEALTHCHECK définie |
| DL011 | 5 | pip install avec --no-cache-dir |
| DL012 | 5 | Pas de pattern curl \| sh |

**Scoring** : commence à 100, soustrait le poids de chaque règle échouée.

### `devkit docker optimize <Dockerfile>`

Suggestions d'amélioration :
- Multi-stage build si absent
- Utilisation d'Alpine ou Distroless
- Merge des instructions RUN
- Optimisation du .dockerignore
- Caching des dépendances (COPY séparé)

### `devkit helm scaffold <app-name>`

Génère un Helm chart opinionated avec bonnes pratiques :

```
<app-name>/
├── Chart.yaml
├── values.yaml           # Replicas, image, service, ingress, resources, HPA, probes
├── .helmignore
└── templates/
    ├── _helpers.tpl       # Labels & selector helpers
    ├── deployment.yaml    # runAsNonRoot, resource limits, probes
    ├── service.yaml       # ClusterIP par défaut
    ├── ingress.yaml       # Conditionnel (.Values.ingress.enabled)
    ├── hpa.yaml           # Autoscaling CPU-based
    ├── serviceaccount.yaml
    └── NOTES.txt          # Instructions post-install
```

**Sécurité par défaut** : runAsNonRoot, runAsUser 1000, resource limits, PDB.

### `devkit k8s check-rbac`

Audit RBAC du cluster Kubernetes :
- ClusterRoleBindings avec cluster-admin
- ServiceAccounts avec permissions wildcard (*)
- Utilisation du ServiceAccount default par des pods
- RoleBindings trop permissifs
- Rôles qui accordent l'accès aux secrets

> Nécessite `k8s.io/client-go` — placeholder dans la version initiale.

### `devkit k8s cost-estimate`

Estimation de coût par namespace basée sur les resource requests :
- Lit les CPU/Memory requests de tous les pods
- Applique les tarifs AWS par défaut (configurable)
- Affiche un tableau : Namespace / CPUs / Memory / $/hour / $/month

> Nécessite `k8s.io/client-go` — placeholder dans la version initiale.

### `devkit ci generate`

Génère des configurations CI/CD :

| Provider | Langages supportés | Fichier généré |
|----------|-------------------|----------------|
| `github` | Go, Python, Node | `.github/workflows/ci.yml` |
| `gitlab` | Go, Python | `.gitlab-ci.yml` |

Templates incluent : tests, linting, build, coverage, artifacts.

## Build & Release

### Build local

```bash
go build -o devkit .
```

### GoReleaser

Configuration `.goreleaser.yml` :
- **OS** : Linux, macOS, Windows
- **Arch** : amd64, arm64
- **CGO** : désactivé (binaires statiques)
- **ldflags** : `-s -w` (strip debug) + injection de version
- **Format** : tar.gz (Linux/macOS), zip (Windows)

### CI/CD

- **ci.yml** : test (race detector + coverage) → lint (golangci-lint) → build (matrix 3 OS x 2 arch)
- **release.yml** : déclenché sur tag `v*` → GoReleaser → GitHub Release avec binaires
