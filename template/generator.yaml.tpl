ignore:
  resource_names:
{{- range $crdName := .CRDNames }}
      - {{ $crdName }}
{{- end }}
{{- if not (eq .ServiceModelName "") }}
model_name: {{ .ServiceModelName }}
{{- end }}
