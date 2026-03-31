# devkit

[![CI](https://github.com/g-holali-david/devkit/actions/workflows/ci.yml/badge.svg)](https://github.com/g-holali-david/devkit/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/g-holali-david/devkit)](https://goreportcard.com/report/github.com/g-holali-david/devkit)

A CLI toolkit for DevOps engineers — lint Dockerfiles, scaffold Helm charts, audit K8s RBAC, estimate cluster costs, and generate CI pipelines.

## Install

```bash
# From source
go install github.com/g-holali-david/devkit@latest

# From release
# Download the binary from https://github.com/g-holali-david/devkit/releases
```

## Commands

### Docker

```bash
# Lint a Dockerfile — quality score with detailed report
devkit docker lint Dockerfile

# Get optimization suggestions
devkit docker optimize Dockerfile
```

Example lint output:
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
  ✗ HEALTHCHECK instruction defined
  ✓ pip install uses --no-cache-dir
  ✓ No use of curl | sh pattern

  ██████████████████░░ 80/100

  10 passed, 2 failed
```

### Helm

```bash
# Scaffold a production-ready Helm chart
devkit helm scaffold my-app --output ./charts
```

Generates:
- `Chart.yaml`, `values.yaml`
- Templates: Deployment, Service, Ingress, HPA, ServiceAccount
- Security defaults: non-root, resource limits, PDB
- Autoscaling, health probes pre-configured

### Kubernetes

```bash
# Audit RBAC permissions
devkit k8s check-rbac

# Estimate cost per namespace
devkit k8s cost-estimate --namespace production
```

### CI Pipeline Generation

```bash
# Generate GitHub Actions workflow
devkit ci generate --provider github --language go --output .

# Generate GitLab CI
devkit ci generate --provider gitlab --language python --output .
```

Supported: `go`, `python`, `node` x `github`, `gitlab`

## Project Structure

```
.
├── cmd/           # Cobra commands
├── pkg/
│   ├── docker/    # Dockerfile lint & optimize
│   ├── helm/      # Chart scaffolding
│   ├── k8s/       # RBAC audit & cost estimation
│   └── ci/        # CI pipeline generation
├── internal/
│   └── output/    # Colored terminal helpers
└── main.go
```

## Development

```bash
# Build
go build -o devkit .

# Test
go test ./...

# Lint
golangci-lint run
```

## Release

Releases are automated via GoReleaser on git tags:

```bash
git tag v0.1.0
git push origin v0.1.0
# GoReleaser builds binaries for Linux/macOS/Windows
```

## License

MIT
