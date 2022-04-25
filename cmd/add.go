/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var (
	appName string

	addCmd = &cobra.Command{
		Use:   "add [module name]",
		Short: "Добавление модуля",
		Long: `Добавление модуля (bctl add) создаст чарт с новым модулем в отдельной директории, 
дополнит чарт с деплоем Argo приложений, добавит нужные values для каждой среды.`,
		Run: func(cmd *cobra.Command, args []string) {
			moduleName, err := addModule(args)
			cobra.CheckErr(err)
			fmt.Printf("Модуль %s добавлен в проект\n", moduleName)
		},
	}
)

func init() {

	addCmd.Flags().StringVar(&projectName, "project-name", "", "fully qualified project name")
	cobra.CheckErr(addCmd.MarkFlagRequired("project-name"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addModule(args []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	wd = fmt.Sprintf("%s/%s-deploy", wd, projectName)
	path := fmt.Sprintf("%s-deploy", projectName)

	if len(args) < 1 {
		return "", errors.New("Не передано название модуля")

	}
	if args[0] != "" {
		appName = fmt.Sprintf("%s", args[0])
	}

	project := &Project{
		Path:         path,
		AbsolutePath: wd,
		ProjectName:  projectName,
		ModuleName:   appName,
	}

	if err := project.AddModule(); err != nil {
		return "", err
	}

	return project.ModuleName, nil
}
