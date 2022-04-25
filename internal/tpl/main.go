package tpl

func ReadmeTemplate() []byte {
	return []byte(`
# armt-deploy

# armt-deploy

install argocd-armt
` + "```" + `bash
helm upgrade --install argocd-armt ./armt-argocd -f armt-argocd/values.yaml --kubeconfig ~/.kube/kubeconfig -n dev-webarm
` + "```")
}

func CiTemplate() []byte {
	return []byte(`
stages:
  - build

include:
  - project: 'tc-center/infra/App/autodeploy'
    ref: master
    file: '.gitlab-ci.yml'
`)
}

func ChartTemplate() []byte {
	return []byte(`
apiVersion: v2
appVersion: 1.0.0 # версия ПРИЛОЖЕНИЯ в чарте
description: A Helm chart for {{ .ProjectName }}-{{ .AppName }}
name: {{ if .AppName }}{{ .AppName }}{{ else }}argocd{{ end }}
type: application
version: 0.1.0 # версия ЧАРТА, обнуляется (0.1.0) с каждым релизом
dependencies:
- name: base
  version: 1.5.1
  repository: "file://./charts/base-1.5.1.tgz"
`)
}

func ArgoValuesTemplate() []byte {
	return []byte(`
base:

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
          - name: {{ .AppName }}
            valuesFiles:
              env: true # env_values/env/env-proj.yaml
              ver: true # module_version.yaml`)
}
