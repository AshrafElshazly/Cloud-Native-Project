{{/*
Return the fully qualified name of the chart.
*/}}
{{- define "product-service.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name }}
{{- end -}}

{{/*
Return the name of the chart.
*/}}
{{- define "product-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end -}}
