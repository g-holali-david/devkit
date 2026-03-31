// Package helm provides Helm chart scaffolding tools.
package helm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/g-holali-david/devkit/internal/output"
)

func Scaffold(appName, outputDir string) error {
	chartDir := filepath.Join(outputDir, appName)
	templatesDir := filepath.Join(chartDir, "templates")

	output.Header("Scaffolding Helm chart: " + appName)

	dirs := []string{chartDir, templatesDir}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("cannot create directory %s: %w", d, err)
		}
	}

	files := map[string]string{
		filepath.Join(chartDir, "Chart.yaml"):               chartYAML(appName),
		filepath.Join(chartDir, "values.yaml"):               valuesYAML(appName),
		filepath.Join(chartDir, ".helmignore"):               helmIgnore(),
		filepath.Join(templatesDir, "_helpers.tpl"):          helpers(appName),
		filepath.Join(templatesDir, "deployment.yaml"):       deployment(),
		filepath.Join(templatesDir, "service.yaml"):          service(),
		filepath.Join(templatesDir, "ingress.yaml"):          ingress(),
		filepath.Join(templatesDir, "hpa.yaml"):              hpa(),
		filepath.Join(templatesDir, "serviceaccount.yaml"):   serviceAccount(),
		filepath.Join(templatesDir, "NOTES.txt"):             notes(appName),
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return fmt.Errorf("cannot write %s: %w", path, err)
		}
		rel, _ := filepath.Rel(outputDir, path)
		output.Pass("Created " + rel)
	}

	fmt.Printf("\n  Chart created at %s\n\n", output.Cyan(chartDir))
	return nil
}

func chartYAML(name string) string {
	return fmt.Sprintf(`apiVersion: v2
name: %s
description: A Helm chart for %s
type: application
version: 0.1.0
appVersion: "1.0.0"
`, name, name)
}

func valuesYAML(name string) string {
	return fmt.Sprintf(`replicaCount: 2

image:
  repository: ghcr.io/g-holali-david/%s
  pullPolicy: IfNotPresent
  tag: ""  # Overridden by CI

serviceAccount:
  create: true
  name: ""

service:
  type: ClusterIP
  port: 80
  targetPort: 8080

ingress:
  enabled: false
  className: nginx
  annotations: {}
  hosts:
    - host: %s.local
      paths:
        - path: /
          pathType: Prefix
  tls: []

resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 256Mi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 75

probes:
  liveness:
    path: /health
    port: http
    initialDelaySeconds: 10
    periodSeconds: 15
  readiness:
    path: /health
    port: http
    initialDelaySeconds: 5
    periodSeconds: 10

podDisruptionBudget:
  enabled: true
  minAvailable: 1

nodeSelector: {}
tolerations: []
affinity: {}
`, name, name)
}

func helmIgnore() string {
	return `.git
.gitignore
*.md
.helmignore
`
}

func helpers(name string) string {
	return fmt.Sprintf(`{{/*
Common labels
*/}}
{{- define "%s.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "%s.selectorLabels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Service account name
*/}}
{{- define "%s.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default .Chart.Name .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
`, name, name, name)
}

func deployment() string {
	return `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Release.Name }}
  labels:
    {{- include "` + "{{ .Chart.Name }}" + `.labels" . | nindent 4 }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "` + "{{ .Chart.Name }}" + `.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "` + "{{ .Chart.Name }}" + `.selectorLabels" . | nindent 8 }}
    spec:
      serviceAccountName: {{ include "` + "{{ .Chart.Name }}" + `.serviceAccountName" . }}
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: {{ .Values.service.targetPort }}
              protocol: TCP
          livenessProbe:
            httpGet:
              path: {{ .Values.probes.liveness.path }}
              port: {{ .Values.probes.liveness.port }}
            initialDelaySeconds: {{ .Values.probes.liveness.initialDelaySeconds }}
            periodSeconds: {{ .Values.probes.liveness.periodSeconds }}
          readinessProbe:
            httpGet:
              path: {{ .Values.probes.readiness.path }}
              port: {{ .Values.probes.readiness.port }}
            initialDelaySeconds: {{ .Values.probes.readiness.initialDelaySeconds }}
            periodSeconds: {{ .Values.probes.readiness.periodSeconds }}
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
`
}

func service() string {
	return `apiVersion: v1
kind: Service
metadata:
  name: {{ .Release.Name }}
  labels:
    {{- include "` + "{{ .Chart.Name }}" + `.labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: {{ .Values.service.targetPort }}
      protocol: TCP
      name: http
  selector:
    {{- include "` + "{{ .Chart.Name }}" + `.selectorLabels" . | nindent 4 }}
`
}

func ingress() string {
	return `{{- if .Values.ingress.enabled -}}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Release.Name }}
  labels:
    {{- include "` + "{{ .Chart.Name }}" + `.labels" . | nindent 4 }}
  {{- with .Values.ingress.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  ingressClassName: {{ .Values.ingress.className }}
  {{- if .Values.ingress.tls }}
  tls:
    {{- range .Values.ingress.tls }}
    - hosts:
        {{- range .hosts }}
        - {{ . | quote }}
        {{- end }}
      secretName: {{ .secretName }}
    {{- end }}
  {{- end }}
  rules:
    {{- range .Values.ingress.hosts }}
    - host: {{ .host | quote }}
      http:
        paths:
          {{- range .paths }}
          - path: {{ .path }}
            pathType: {{ .pathType }}
            backend:
              service:
                name: {{ $.Release.Name }}
                port:
                  number: {{ $.Values.service.port }}
          {{- end }}
    {{- end }}
{{- end }}
`
}

func hpa() string {
	return `{{- if .Values.autoscaling.enabled }}
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: {{ .Release.Name }}
  labels:
    {{- include "` + "{{ .Chart.Name }}" + `.labels" . | nindent 4 }}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ .Release.Name }}
  minReplicas: {{ .Values.autoscaling.minReplicas }}
  maxReplicas: {{ .Values.autoscaling.maxReplicas }}
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: {{ .Values.autoscaling.targetCPUUtilizationPercentage }}
{{- end }}
`
}

func serviceAccount() string {
	return `{{- if .Values.serviceAccount.create -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "` + "{{ .Chart.Name }}" + `.serviceAccountName" . }}
  labels:
    {{- include "` + "{{ .Chart.Name }}" + `.labels" . | nindent 4 }}
{{- end }}
`
}

func notes(name string) string {
	return strings.ReplaceAll(`1. Get the application URL:
{{- if .Values.ingress.enabled }}
  http{{ if $.Values.ingress.tls }}s{{ end }}://{{ (index .Values.ingress.hosts 0).host }}
{{- else }}
  kubectl port-forward svc/{{ .Release.Name }} 8080:{{ .Values.service.port }}
  Then visit http://localhost:8080
{{- end }}
`, "\t", "  ")
}
