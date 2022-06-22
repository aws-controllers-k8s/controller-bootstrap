ignore:
  resource_names:
{{- range $crdName := .CRDNames }}
      - {{ $crdName }}
{{- end }}
{{ $serviceModelName := .ServiceModelName }}
{{- if not (eq $serviceModelName "") -}}
    {{ $serviceModelName = .ServiceModelName }}
model_name: {{ $serviceModelName }}
{{- end }}
