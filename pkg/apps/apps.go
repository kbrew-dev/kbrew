package apps

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/kbrew-dev/kbrew/pkg/apps/helm"
	"github.com/kbrew-dev/kbrew/pkg/apps/raw"
	"github.com/kbrew-dev/kbrew/pkg/config"
)

// Method defines operation performed on the apps
type Method string

const (
	// Install method to install the app
	Install Method = "install"
	// Uninstall method to uninstall the app
	Uninstall Method = "uninstall"
)

// App represents a K8s applications than can be managed with kbrew recipes
type App interface {
	Install(ctx context.Context, name, namespace string, version string, opt map[string]string) error
	Uninstall(ctx context.Context, name, namespace string) error
	Search(ctx context.Context, name string) (string, error)
}

// Run fetches recipe from registry for the app and performs given operation
func Run(ctx context.Context, m Method, appName, namespace, appConfigPath string) error {
	c, err := config.New(appConfigPath)
	if err != nil {
		return err
	}
	var app App

	switch c.App.Repository.Type {
	case config.Helm:
		app = helm.New(c.App)
	case config.Raw:
		app, err = raw.New(c.App)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("Unsupported app type %s", c.App.Repository.Type))
	}

	// Check if entry exists in config
	if c.App.Name != appName {
		// Check if app exists in repo
		if _, err := app.Search(ctx, appName); err != nil {
			return err
		}
	}

	// Override if default namespace is set
	if c.App.Namespace != "" {
		namespace = c.App.Namespace
	}
	if c.App.Namespace == "-" {
		namespace = ""
	}

	switch m {
	case Install:
		// Run preinstall
		for _, phase := range c.App.PreInstall {
			for _, a := range phase.Apps {
				if err := Run(ctx, m, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
					return err
				}
			}
			for _, a := range phase.Steps {
				if err := execCommand(a); err != nil {
					return err
				}
			}
		}
		// Run install
		if err := app.Install(ctx, appName, namespace, c.App.Version, nil); err != nil {
			return err
		}
		// Run postinstall
		for _, phase := range c.App.PostInstall {
			for _, a := range phase.Apps {
				if err := Run(ctx, m, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
					return err
				}
			}
			for _, a := range phase.Steps {
				if err := execCommand(a); err != nil {
					return err
				}
			}
		}
	case Uninstall:
		return app.Uninstall(ctx, appName, namespace)
	default:
		return errors.New(fmt.Sprintf("Unsupported method %s", m))
	}
	return nil
}

func execCommand(cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
