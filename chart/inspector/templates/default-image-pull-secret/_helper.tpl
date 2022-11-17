{{/*
Return the Default Image Pull Secret DockerConfigJson
*/}}
{{- define "inspector.defaultImagePullSecret.dockerconfigjson" -}}
  {{- if .Values.defaultImagePullSecret.enabled -}}
    {{- $registry := .Values.global.image.registry -}}
    {{- with .Values.defaultImagePullSecret -}}
      {{- printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"email\":\"%s\",\"auth\":\"%s\"}}}" $registry .username .password .email (printf "%s:%s" .username .password | b64enc) | b64enc }}
    {{- end -}}
  {{- end -}}
{{- end -}}

