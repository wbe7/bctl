/*
Copyright © 2022 Mikhail Tikhonov mikhail.tikhonov@stoloto.ru

*/

package cmd

import (
	"fmt"
	"os"
	"text/template"

	"git.tccenter.ru/tc-center/infra/App/base/bctl/internal/tpl"
	"github.com/spf13/cobra"
)

// Project contains name, license and paths to projects.
type Project struct {
	// v2
	Path          string
	ProjectName   string
	AbsolutePath  string
	ModuleName    string
	ModuleImage   string
	ModuleVersion string
	ModulePort    string
	IngressClass  string
	ChartVersion  string
}

func (p *Project) Create() error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			return err
		}
	}

	// create .gitlab-ci.yml
	ciFile, err := os.Create(fmt.Sprintf("%s/.gitlab-ci.yml", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer ciFile.Close()

	ciTemplate := template.Must(template.New("ci").Parse(string(tpl.CiTemplate())))
	err = ciTemplate.Execute(ciFile, p)
	if err != nil {
		return err
	}

	// create README.md
	readmeFile, err := os.Create(fmt.Sprintf("%s/README.md", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer readmeFile.Close()

	readmeTemplate := template.Must(template.New("readme").Parse(string(tpl.ReadmeTemplate())))
	err = readmeTemplate.Execute(readmeFile, p)
	if err != nil {
		return err
	}

	// create {{ .ProjectName }}-argocd/Chart.yaml
	createDir(p, fmt.Sprintf("%s-argocd", p.ProjectName))

	chartFile, err := os.Create(fmt.Sprintf("%s/%s-argocd/Chart.yaml", p.AbsolutePath, p.ProjectName))
	if err != nil {
		return err
	}
	defer chartFile.Close()

	chartTemplate := template.Must(template.New("chart").Parse(string(tpl.ChartTemplate())))
	err = chartTemplate.Execute(chartFile, p)
	if err != nil {
		return err
	}

	// create {{ .ProjectName }}-argocd/values.yaml
	valuesFile, err := os.Create(fmt.Sprintf("%s/%s-argocd/values.yaml", p.AbsolutePath, p.ProjectName))
	if err != nil {
		return err
	}
	defer valuesFile.Close()

	valuesTemplate := template.Must(template.New("values").Parse(string(tpl.ArgoValuesTemplate())))
	err = valuesTemplate.Execute(valuesFile, p)
	if err != nil {
		return err
	}

	// create env-values/{ENV}
	createDir(p, "env-values")

	createDir(p, "env-values/dev")

	createDir(p, "env-values/tfi")

	createDir(p, "env-values/tifa")

	createDir(p, "env-values/tli")

	createDir(p, "env-values/prod")

	return nil
}

func (p *Project) AddModule() error {

	// create {{ .ProjectName }}-argocd/values.yaml
	createDir(p, fmt.Sprintf("%s-argocd", p.ProjectName))

	argoValuesFile, err := os.OpenFile(fmt.Sprintf("%s/%s-argocd/values.yaml", p.AbsolutePath, p.ProjectName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer argoValuesFile.Close()

	valuesTemplate := template.Must(template.New("values").Parse(string(tpl.ArgoModuleValuesTemplate())))
	err = valuesTemplate.Execute(argoValuesFile, p)
	if err != nil {
		return err
	}

	// create {{ .ProjectName }}-{{ .ModuleName}}/Chart.yaml
	createDir(p, fmt.Sprintf("%s-%s", p.ProjectName, p.ModuleName))

	if _, err := os.Stat(fmt.Sprintf("%s/%s-%s", p.AbsolutePath, p.ProjectName, p.ModuleName)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s-%s", p.AbsolutePath, p.ProjectName, p.ModuleName), 0751))
	}
	chartFile, err := os.Create(fmt.Sprintf("%s/%s-%s/Chart.yaml", p.AbsolutePath, p.ProjectName, p.ModuleName))
	if err != nil {
		return err
	}
	defer chartFile.Close()

	chartTemplate := template.Must(template.New("chart").Parse(string(tpl.ChartTemplate())))
	err = chartTemplate.Execute(chartFile, p)
	if err != nil {
		return err
	}

	// create {{ .ProjectName }}-{{ .ModuleName}}/values.yaml
	valuesFile, err := os.Create(fmt.Sprintf("%s/%s-%s/values.yaml", p.AbsolutePath, p.ProjectName, p.ModuleName))
	if err != nil {
		return err
	}
	defer valuesFile.Close()

	valuesTemplate = template.Must(template.New("chart").Parse(string(tpl.ModuleValuesTemplate())))
	err = valuesTemplate.Execute(valuesFile, p)
	if err != nil {
		return err
	}

	// create {{ .ProjectName }}-{{ .ModuleName}}/module_version.yaml
	moduleVersionFile, err := os.Create(fmt.Sprintf("%s/%s-%s/module_version.yaml", p.AbsolutePath, p.ProjectName, p.ModuleName))
	if err != nil {
		return err
	}
	defer moduleVersionFile.Close()

	valuesTemplate = template.Must(template.New("moduleVersion").Parse(string(tpl.ModuleVersionTemplate())))
	err = valuesTemplate.Execute(moduleVersionFile, p)
	if err != nil {
		return err
	}

	// create values env-values
	err = createEnvValuesFile(p, "dev")
	if err != nil {
		return err
	}

	err = createEnvValuesFile(p, "tfi")
	if err != nil {
		return err
	}

	err = createEnvValuesFile(p, "tifa")
	if err != nil {
		return err
	}

	err = createEnvValuesFile(p, "tli")
	if err != nil {
		return err
	}

	err = createEnvValuesFile(p, "prod")
	if err != nil {
		return err
	}

	return nil
}

func createEnvValuesFile(p *Project, env string) error {
	envValuesFile, err := os.Create(fmt.Sprintf("%s/env-values/"+env+"/"+env+"-%s-%s.yaml", p.AbsolutePath, p.ProjectName, p.ModuleName))
	if err != nil {
		return err
	}
	defer envValuesFile.Close()
	return nil
}

func createDir(p *Project, dir string) {
	if _, err := os.Stat(fmt.Sprintf("%s/"+dir, p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/"+dir, p.AbsolutePath), 0751))
	}
}
