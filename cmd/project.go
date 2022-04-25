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
	AppName      string
}

type Module struct {
	ModuleName   string
	ModuleParent string
	*Project
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
	defer chartFile.Close()

	valuesTemplate := template.Must(template.New("values").Parse(string(tpl.ArgoValuesTemplate())))
	err = valuesTemplate.Execute(valuesFile, p)
	if err != nil {
		return err
	}
	return nil
}

//func (c *Module) Create() error {
//	cmdFile, err := os.Create(fmt.Sprintf("%s/cmd/%s.go", c.AbsolutePath, c.ModuleName))
//	if err != nil {
//		return err
//	}
//	defer cmdFile.Close()
//
//	commandTemplate := template.Must(template.New("sub").Parse(string(tpl.AddCommandTemplate())))
//	err = commandTemplate.Execute(cmdFile, c)
//	if err != nil {
//		return err
//	}
//	return nil
//}
