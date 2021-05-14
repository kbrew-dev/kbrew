package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kbrew-dev/kbrew/pkg/apps"
	"github.com/kbrew-dev/kbrew/pkg/config"
	"github.com/kbrew-dev/kbrew/pkg/registry"
	"github.com/kbrew-dev/kbrew/pkg/update"
	"github.com/kbrew-dev/kbrew/pkg/version"
)

const defaultTimeout = "15m0s"

var (
	configFile string
	namespace  string
	timeout    string

	rootCmd = &cobra.Command{
		Use:   "kbrew",
		Short: "Homebrew for your Kubernetes applications",
		Long:  `TODO: Long description`,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Long(context.Background()))
		},
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

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update kbrew and recipe registries",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Upgrade kbrew
			if err := update.CheckRelease(context.Background()); err != nil {
				return err
			}
			// Update kbrew registries
			reg, err := registry.New(config.ConfigDir)
			if err != nil {
				return err
			}
			return reg.Update()
		},
	}

	analyticsCmd = &cobra.Command{
		Use:   "analytics [on|off|status]",
		Short: "Manage analytics setting",
		RunE: func(cmd *cobra.Command, args []string) error {
			return manageAnalytics(args)
		},
	}
)

func init() {
	cobra.OnInitialize(config.InitConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.kbrew.yaml)")
	rootCmd.PersistentFlags().StringVarP(&config.ConfigDir, "config-dir", "", "", "config dir (default is $HOME/.kbrew)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(analyticsCmd)

	installCmd.PersistentFlags().StringVarP(&timeout, "timeout", "t", "", "time to wait for app components to be in a ready state (default 15m0s)")
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
		errors.New("No app name provided")
	}
	return nil
}

func manageApp(m apps.Method, args []string) error {
	ctx := context.Background()
	if timeout == "" {
		timeout = defaultTimeout
	}
	timeoutDur, err := time.ParseDuration(timeout)
	if err != nil {
		return err
	}
	for _, a := range args {
		reg, err := registry.New(config.ConfigDir)
		if err != nil {
			return err
		}
		configFile, err := reg.FetchRecipe(strings.ToLower(a))
		if err != nil {
			return err
		}
		ctxTimeout, cancel := context.WithTimeout(ctx, timeoutDur)
		defer cancel()
		if err := apps.Run(ctxTimeout, m, strings.ToLower(a), namespace, configFile); err != nil {
			return err
		}
	}
	return nil
}

func manageAnalytics(args []string) error {
	if len(args) == 0 {
		return errors.New("Missing subcommand")
	}
	switch args[0] {
	case "on":
		viper.Set(config.AnalyticsEnabled, true)
		return viper.WriteConfig()
	case "off":
		viper.Set(config.AnalyticsEnabled, false)
		return viper.WriteConfig()
	case "status":
		kc, err := config.NewKbrew()
		if err != nil {
			return err
		}
		fmt.Println("Analytics enabled:", kc.AnalyticsEnabled)
	default:
		return errors.New("Invalid subcommand")
	}
	return nil
}
