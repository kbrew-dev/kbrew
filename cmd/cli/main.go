package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/vishal-biyani/kbrew/pkg/apps"
	"github.com/vishal-biyani/kbrew/pkg/apps/raw"
)

type method string

const (
	create method = "create"
	delete method = "delete"
)

var (
	rootCmd = &cobra.Command{
		Use:   "kbrew",
		Short: "Homebrew for your Kubernetes applications",
		Long:  `TODO: Long description`,
	}

	installCmd = &cobra.Command{
		Use:   "install [NAME]",
		Short: "Install application",
		RunE: func(cmd *cobra.Command, args []string) error {
			return manageApp(create, args)
		},
	}

	removeCmd = &cobra.Command{
		Use:   "remove [NAME]",
		Short: "Remove application",
		RunE: func(cmd *cobra.Command, args []string) error {
			return manageApp(delete, args)
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

func manageApp(m method, args []string) error {
	var app apps.App
	var err error
	for _, a := range args {
		app, err = raw.New(a)
		if err != nil {
			// TODO: Check other app types
			return err
		}
		data, err := app.Manifest(context.Background(), nil)
		if err != nil {
			return err
		}
		//fmt.Printf("Data:: \n%s\n", data)
		if err := executeCommand(m, data); err != nil {
			return err
		}
	}
	return nil
}

func search(args []string) error {
	ctx := context.Background()
	// List raw apps
	rawApps, err := raw.List(ctx)
	if err != nil {
		return err
	}
	//TODO: Support for other app types
	if len(args) == 0 {
		printList(rawApps)
		return nil
	}
	result := []string{}
	for _, a := range rawApps {
		if strings.HasPrefix(a, args[0]) {
			result = append(result, a)
		}
	}
	printList(result)
	return nil
}

func printList(list []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "NAME\tVERSION")
	for _, l := range list {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s", l, "NA"))
	}
	w.Flush()

}

func executeCommand(m method, data []byte) error {
	// Generate code
	c := exec.Command("kubectl", string(m), "-f", "-")
	c.Stdin = strings.NewReader(string(data))
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
