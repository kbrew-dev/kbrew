// Copyright 2021 The kbrew Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/kbrew-dev/kbrew/pkg/apps"
	"github.com/kbrew-dev/kbrew/pkg/config"
	"github.com/kbrew-dev/kbrew/pkg/log"
	"github.com/kbrew-dev/kbrew/pkg/registry"
	"github.com/kbrew-dev/kbrew/pkg/update"
	"github.com/kbrew-dev/kbrew/pkg/version"
)

const defaultTimeout = "15m0s"

var (
	configFile string
	namespace  string
	timeout    string
	debug      bool

	rootCmd = &cobra.Command{
		Use:           "kbrew",
		Short:         "A Kubernetes tool that aims to make installing any complex stack in any cloud possible with one step.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Long())
			release, err := update.IsAvailable(context.Background())
			if err != nil {
				fmt.Printf("Error getting latest version of kbrew from GiThub: %s", err)
			}
			if release != "" {
				fmt.Printf("There is a new version of kbrew available: %s, please run 'kbrew update' command to upgrade.\n", release)
			}
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

	completionCmd = &cobra.Command{
		Use:       "completion [SHELL]",
		Short:     "Output shell completion code for the specified shell",
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			switch args[0] {
			case "bash":
				err = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				err = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				err = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return err
		},
	}

	infoCmd = &cobra.Command{
		Use:   "info [NAME]",
		Short: "Describe application",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, err := registry.New(config.ConfigDir)
			if err != nil {
				return err
			}
			s, err := reg.Info(args[0])
			if err != nil {
				return err
			}
			fmt.Println(s)
			return nil
		},
	}

	argsCmd = &cobra.Command{
		Use:   "args [NAME]",
		Short: "Get arguments for an application",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, err := registry.New(config.ConfigDir)
			if err != nil {
				return err
			}
			appArgs, err := reg.Args(args[0])
			if err != nil {
				return err
			}

			bytes, err := yaml.Marshal(appArgs)
			if err != nil {
				return err
			}
			fmt.Println(string(bytes))
			return nil
		},
	}
)

func init() {
	cobra.OnInitialize(config.InitConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.kbrew.yaml)")
	rootCmd.PersistentFlags().StringVarP(&config.ConfigDir, "config-dir", "", "", "config dir (default is $HOME/.kbrew)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "enable debug logs")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(analyticsCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(infoCmd)

	infoCmd.AddCommand(argsCmd)

	installCmd.PersistentFlags().StringVarP(&timeout, "timeout", "t", "", "time to wait for app components to be in a ready state (default 15m0s)")
}

func main() {
	Execute()
}

// Execute executes the main command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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
		logger := log.NewLogger(debug)
		runner := apps.NewAppRunner(m, logger, log.NewStatus(logger))
		c, err := config.NewApp(strings.ToLower(a), configFile)
		if err != nil {
			return err
		}
		printDetails(logger, strings.ToLower(a), m, c)
		ctxTimeout, cancel := context.WithTimeout(ctx, timeoutDur)
		defer cancel()
		if err := runner.Run(ctxTimeout, strings.ToLower(a), namespace, configFile); err != nil {
			return err
		}
	}
	return nil
}

func printDetails(log *log.Logger, appName string, m apps.Method, c *config.AppConfig) {
	switch m {
	case apps.Install:
		log.Infof("🚀 Installing %s app...", appName)
		log.InfoMap("Version", c.App.Version)
		log.InfoMap("Pre-install dependencies", "")
		for _, pre := range c.App.PreInstall {
			for _, app := range pre.Apps {
				log.Infof(" - %s", app)
			}
		}
		log.InfoMap("Post-install dependencies", "")
		for _, post := range c.App.PostInstall {
			for _, app := range post.Apps {
				log.Infof(" - %s", app)
			}
		}
		log.Info("---")
	case apps.Uninstall:
		log.Infof("🧹 Uninstalling %s app and its dependencies...", appName)
		log.InfoMap("Dependencies", "")
		for _, pre := range c.App.PreInstall {
			for _, app := range pre.Apps {
				log.Infof(" - %s", app)
			}
		}
		for _, post := range c.App.PostInstall {
			for _, app := range post.Apps {
				log.Infof(" - %s", app)
			}
		}
		log.Info("---")
	}

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
