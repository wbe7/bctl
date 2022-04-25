package tpl

func ReadmeTemplate() []byte {
	return []byte(`# {{ .ProjectName }}

## {{ .Path }}

Install argocd-{{ .ProjectName }}
` + "```" + `bash
helm upgrade --install argocd-{{ .ProjectName }} ./{{ .ProjectName }}-argocd -f {{ .ProjectName }}-argocd/values.yaml --kubeconfig ~/.kube/kubeconfig -n dev-{{ .ProjectName }}
` + "```")
}

func CiTemplate() []byte {
	return []byte(`stages:
  - build

include:
  - project: 'tc-center/infra/App/base/autodeploy'
    ref: master
    file: '.gitlab-ci.yml'
`)
}

func ChartTemplate() []byte {
	return []byte(`apiVersion: v2
appVersion: 1.0.0 # версия ПРИЛОЖЕНИЯ в чарте
description: A Helm chart for {{ .ProjectName }}-{{ if .ModuleName }}{{ .ModuleName }}{{ else }}argocd{{ end }}
name: {{ if .ModuleName }}{{ .ModuleName }}{{ else }}argocd{{ end }}
type: application
version: 0.1.0 # версия ЧАРТА, обнуляется (0.1.0) с каждым релизом
dependencies:
- name: base
  version: 1.5.1
  repository: "file://./charts/base-1.5.1.tgz"
`)
}

func ArgoValuesTemplate() []byte {
	return []byte(`base:

  chartName: "{{ .ProjectName }}"

  argocd:
	#TODO
    telegram: "CHANGEME" # id телеграм чата
    repo: # список гит реп
      - repoLink: "git@git.tccenter.ru:tc-center/infra/App/{{ .ProjectName }}-deploy.git" # ссылка на репу
        path: {{ .ProjectName }}
        repoApps: # список приложений`)
}

func ArgoModuleValuesTemplate() []byte {
	return []byte(`
          - name: {{ .ModuleName }}
            valuesFiles:
              env: true # env_values/env/env-proj.yaml
              ver: true # module_version.yaml`)
}

func ModuleValuesTemplate() []byte {
	return []byte(`base:

  chartName: "{{ .ProjectName }}"
  hashicorpv1: true

  deployment:
    replicas: 1
`)
}
