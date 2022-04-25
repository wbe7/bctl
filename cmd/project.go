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
	Path         string
	ProjectName  string
	AbsolutePath string
	ModuleName   string
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
	if _, err = os.Stat(fmt.Sprintf("%s/%s-argocd", p.AbsolutePath, p.ProjectName)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s-argocd", p.AbsolutePath, p.ProjectName), 0751))
	}
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
	if _, err = os.Stat(fmt.Sprintf("%s/%s-argocd", p.AbsolutePath, p.ProjectName)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s-argocd", p.AbsolutePath, p.ProjectName), 0751))
	}
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
	return nil
}

func (p *Project) AddModule() error {

	// create {{ .ProjectName }}-argocd/values.yaml
	if _, err := os.Stat(fmt.Sprintf("%s/%s-argocd", p.AbsolutePath, p.ProjectName)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s-argocd", p.AbsolutePath, p.ProjectName), 0751))
	}
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
	if _, err := os.Stat(fmt.Sprintf("%s/%s-%s", p.AbsolutePath, p.ProjectName, p.ModuleName)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s-%s", p.AbsolutePath, p.ProjectName, p.ModuleName), 0751))
	}
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

	return nil
}
