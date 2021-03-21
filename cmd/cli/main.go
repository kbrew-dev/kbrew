package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kbrew-dev/kbrew/pkg/apps"
	"github.com/kbrew-dev/kbrew/pkg/config"
	"github.com/kbrew-dev/kbrew/pkg/registry"
)

var (
	configFile string
	namespace  string

	rootCmd = &cobra.Command{
		Use:   "kbrew",
		Short: "Homebrew for your Kubernetes applications",
		Long:  `TODO: Long description`,
	}

	installCmd = &cobra.Command{
		Use:   "install [NAME]",
		Short: "Install application",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return manageApp(apps.Install, args)
		},
	}

	removeCmd = &cobra.Command{
		Use:   "remove [NAME]",
		Short: "Remove application",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return manageApp(apps.Uninstall, args)
		},
	}

	searchCmd = &cobra.Command{
		Use:   "search [NAME]",
		Short: "Search application",
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := ""
			if len(args) != 0 {
				appName = args[0]
			}
			reg, err := registry.New(config.ConfigDir)
			if err != nil {
				return err
			}
			appList, err := reg.Search(appName, false)
			if err != nil {
				return err
			}
			if len(appList) == 0 {
				fmt.Printf("No recipe found for %s.\n", appName)
				return nil
			}
			fmt.Println("Available recipes:")
			for _, app := range appList {
				fmt.Println(app.Name)
			}
			return nil
		},
	}
)

func init() {
	cobra.OnInitialize(config.InitConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.kbrew.yaml)")
	rootCmd.PersistentFlags().StringVarP(&config.ConfigDir, "config-dir", "", "", "config dir (default is $HOME/.kbrew)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace")

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(searchCmd)
}

func main() {
	Execute()
}

// Execute executes the main command
func Execute() error {
	return rootCmd.Execute()
}

func checkArgs(args []string) error {
	if len(args) == 0 {
		errors.New("No app name provided.")
	}
	return nil
}

func manageApp(m apps.Method, args []string) error {
	ctx := context.Background()
	for _, a := range args {
		reg, err := registry.New(config.ConfigDir)
		if err != nil {
			return err
		}
		configFile, err := reg.FetchRecipe(strings.ToLower(a))
		if err != nil {
			return err
		}
		if err := apps.Run(ctx, m, strings.ToLower(a), namespace, configFile); err != nil {
			return err
		}
	}
	return nil
}
