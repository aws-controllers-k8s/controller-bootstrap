ignore:
  resource_names:
{{- range $crdName := .CRDNames }}
      - {{ $crdName }}
{{- end }}
{{- if not (eq .ServiceModelName "") }}
sdk_names:
  model_name: {{ .ServiceModelName }}
{{- end }}
