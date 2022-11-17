{{/* vim: set filetype=mustache: */}}

{{/*
Return the proper virgo storage class name
*/}}
{{- define "inspector.controller" -}}
  {{- printf "inspector-controller" -}}
{{- end -}}

{{/*
Allow the release namespace to be overridden for multi-namespace deployments in combined charts.
*/}}
{{- define "common.names.namespace" -}}
{{- if .Values.namespaceOverride -}}
{{- .Values.namespaceOverride -}}
{{- else -}}
{{- .Release.Namespace -}}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for deployment.
*/}}
{{- define "common.capabilities.deployment.apiVersion" -}}
{{- if semverCompare "<1.14-0" (include "common.capabilities.kubeVersion" .) -}}
{{- print "extensions/v1beta1" -}}
{{- else -}}
{{- print "apps/v1" -}}
{{- end -}}
{{- end -}}

{{/*
Return the target Kubernetes version
*/}}
{{- define "common.capabilities.kubeVersion" -}}
{{- if .Values.global }}
    {{- if .Values.global.kubeVersion }}
    {{- .Values.global.kubeVersion -}}
    {{- else }}
    {{- default .Capabilities.KubeVersion.Version .Values.kubeVersion -}}
    {{- end -}}
{{- else }}
{{- default .Capabilities.KubeVersion.Version .Values.kubeVersion -}}
{{- end -}}
{{- end -}}

{{- define "common.images.image" -}}
{{- $registryName := .imageRoot.registry -}}
{{- $projectName := .imageRoot.project -}}
{{- $repositoryName := .imageRoot.repository -}}
{{- $separator := ":" -}}
{{- $termination := .imageRoot.tag | toString -}}
{{- if .global }}
  {{- if .global.image.registry }}
    {{- $registryName = .global.image.registry -}}
  {{- end -}}
  {{- if .global.image.project }}
    {{- $projectName = .global.image.project -}}
  {{- end -}}
  {{- if .global.image.tag }}
    {{- $termination = .global.image.tag -}}
  {{- end -}}
{{- end -}}
{{- if .imageRoot.digest }}
    {{- $separator = "@" -}}
    {{- $termination = .imageRoot.digest | toString -}}
{{- end -}}
{{- printf "%s/%s/%s%s%s" $registryName $projectName $repositoryName $separator $termination -}}
{{- end -}}

{{/*
Return the proper virgo storage class provisioner image name
*/}}
{{- define "inspector.controller.image" -}}
  {{ include "common.images.image" (dict "imageRoot" .Values.controller.image "global" .Values.global) }}
{{- end -}}


{{/*
Return the appropriate apiVersion for RBAC resources.
*/}}
{{- define "common.capabilities.rbac.apiVersion" -}}
{{- if semverCompare "<1.17-0" (include "common.capabilities.kubeVersion" .) -}}
{{- print "rbac.authorization.k8s.io/v1beta1" -}}
{{- else -}}
{{- print "rbac.authorization.k8s.io/v1" -}}
{{- end -}}
{{- end -}}

{{/*
Return the Default Image Pull Secret Name
*/}}
{{- define "inspector.defaultImagePullSecret" -}}
  {{- printf "inspector-default-image-pull-secret" -}}
{{- end -}}