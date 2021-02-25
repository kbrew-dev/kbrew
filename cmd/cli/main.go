package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/vishal-biyani/kbrew/pkg/apps"
	"github.com/vishal-biyani/kbrew/pkg/apps/helm"
	"github.com/vishal-biyani/kbrew/pkg/apps/raw"
	"github.com/vishal-biyani/kbrew/pkg/config"
)

type method string

const (
	install   method = "create"
	uninstall method = "uninstall"
)

var (
	configFile string
	namespace  string
	version    string

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
			return manageApp(install, args)
		},
	}

	removeCmd = &cobra.Command{
		Use:   "remove [NAME]",
		Short: "Remove application",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return manageApp(uninstall, args)
		},
	}

	searchCmd = &cobra.Command{
		Use:   "search [NAME]",
		Short: "Search application",
		RunE: func(cmd *cobra.Command, args []string) error {
			return search(args)
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.kbrew.yaml)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace")
	installCmd.Flags().StringVarP(&version, "version", "v", "", "App version to be installed")

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

func manageApp(m method, args []string) error {
	ctx := context.Background()
	c, err := config.New(configFile)
	if err != nil {
		return nil
	}

	for _, a := range args {
		var app apps.App
		installApp := strings.ToLower(a)

		switch c.App.Repository.Type {
		case config.Helm:
			app = helm.New(c.App, namespace)
		case config.Raw:
			app, err = raw.New(c.App, namespace)
			if err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("Unsupported app type %s", c.App.Repository.Type))
		}

		if version == "" && c.App.Name == installApp {
			version = c.App.Version
		}

		// Check if entry exists in config
		if c.App.Name != installApp {
			// Check if app exists in repo
			if _, err := app.Search(ctx, installApp); err != nil {
				continue
			}
		}

		switch m {
		case install:
			return app.Install(ctx, installApp, version, nil)
		case uninstall:
			return app.Uninstall(ctx, installApp)
		default:
			return errors.New(fmt.Sprintf("Unsupported method %s", m))
		}

	}
	return nil
}

func search(args []string) error {
	ctx := context.Background()
	c, err := config.New(configFile)
	if err != nil {
		return err
	}

	var app apps.App
	switch c.App.Repository.Type {
	case config.Helm:
		app = helm.New(c.App, namespace)
	case config.Raw:
		app, err = raw.New(c.App, namespace)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("Unsupported app type %s", c.App.Repository.Type))
	}

	if len(args) == 0 {
		out, err := app.Search(ctx, "")
		fmt.Print(string(out))
		return err
	}
	out, err := app.Search(ctx, args[0])
	fmt.Print(string(out))
	return err
}

func initConfig() {
	if configFile != "" {
		return
	}
	// Find home directory.
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	// Generate default config file path
	configFile = filepath.Join(home, ".kbrew.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create file with default config
		c := []byte("apiVersion: v1\nkind: kbrew\napps:\n")
		err := ioutil.WriteFile(configFile, c, 0644)
		cobra.CheckErr(err)
	}
}
