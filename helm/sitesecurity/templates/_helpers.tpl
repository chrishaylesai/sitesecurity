{{- define "sitesecurity.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "sitesecurity.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- define "sitesecurity.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}

{{- define "sitesecurity.dbHost" -}}
{{ include "sitesecurity.fullname" . }}-db
{{- end }}

{{- define "sitesecurity.authHost" -}}
{{ include "sitesecurity.fullname" . }}-auth
{{- end }}

{{- define "sitesecurity.apiHost" -}}
{{ include "sitesecurity.fullname" . }}-api
{{- end }}

{{- define "sitesecurity.frontendHost" -}}
{{ include "sitesecurity.fullname" . }}-frontend
{{- end }}
