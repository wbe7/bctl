/*
Copyright © 2022 Mikhail Tikhonov mikhail.tikhonov@stoloto.ru

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// initCmd represents the init command
var (
	projectName string

	initCmd = &cobra.Command{
		Use:   "init [name]",
		Short: "Инициализация проекта",
		Long: `Инициализация (bctl init) создаст новый проект с минимальным функционалом,
шаблоном для деплоя Argo приложений и автоматизацией автодеплоя.
`,
		Run: func(cmd *cobra.Command, args []string) {
			projectPath, err := initializeProject()
			cobra.CheckErr(err)
			fmt.Printf("Проект деплоя инициализирован в\n%s\n", projectPath)
		},
	}
)

func init() {
	initCmd.Flags().StringVar(&projectName, "project-name", "", "Имя проекта (обязательное поле)")
	cobra.CheckErr(initCmd.MarkFlagRequired("project-name"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initializeProject() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	wd = fmt.Sprintf("%s/%s-deploy", wd, projectName)
	path := fmt.Sprintf("%s-deploy", projectName)

	project := &Project{
		Path:         path,
		AbsolutePath: wd,
		ProjectName:  projectName,
	}

	if err := project.Create(); err != nil {
		return "", err
	}

	return project.AbsolutePath, nil
}
