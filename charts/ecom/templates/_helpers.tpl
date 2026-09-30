{{- define "ecom.fullname" -}}
{{ .Release.Name }}
{{- end -}}

{{- define "ecom.labels" -}}
app.kubernetes.io/part-of: ecom
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/* Bitnami postgresql chart's service is named "<release>-postgresql" */}}
{{- define "ecom.databaseURL" -}}
postgresql://{{ .Values.postgresql.auth.username }}:{{ .Values.postgresql.auth.password }}@{{ .Release.Name }}-postgresql:5432/{{ .Values.postgresql.auth.database }}
{{- end -}}

{{/* Bitnami redis chart's standalone service is named "<release>-redis-master" */}}
{{- define "ecom.redisAddr" -}}
{{ .Release.Name }}-redis-master:6379
{{- end -}}
