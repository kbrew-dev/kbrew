package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

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
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.kbrew.yaml)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "namespace")

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

	// TODO: Switch to config based registration
	for _, a := range args {
		var app apps.App
		// Check if entry exists in config
		// TODO: Create a map during init
		if c.App.Name != strings.ToLower(a) {
			continue
		}

		switch c.App.Type {
		case config.Helm:
			app = helm.New(c.App, namespace)
		case config.Raw:
			app = raw.New(c.App, namespace)
		default:
			return errors.New(fmt.Sprintf("Unsupported app type %s", c.App.Type))
		}

		switch m {
		case install:
			return app.Install(ctx, nil)
		case uninstall:
			return app.Uninstall(ctx)
		default:
			return errors.New(fmt.Sprintf("Unsupported method %s", m))
		}

	}
	return nil
}

func search(args []string) error {

	// List from config
	c, err := config.New(configFile)
	if err != nil {
		return err
	}

	//TODO: Support for other app types
	if len(args) == 0 || strings.HasPrefix(c.App.Name, args[0]) {
		printList(c.App)
		return nil
	}
	return nil
}

func printList(app config.App) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "NAME\tVERSION\tTYPE")
	fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", app.Name, app.Version, app.Type))
	w.Flush()

}

func executeCommand(m method, data []byte) error {
	c := exec.Command("kubectl", string(m), "-f", "-")
	c.Stdin = strings.NewReader(string(data))
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
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
