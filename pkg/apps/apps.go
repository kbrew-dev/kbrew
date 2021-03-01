package apps

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/infracloudio/kbrew/pkg/apps/helm"
	"github.com/infracloudio/kbrew/pkg/apps/raw"
	"github.com/infracloudio/kbrew/pkg/config"
)

type Method string

const (
	Install   Method = "install"
	Uninstall Method = "uninstall"
)

type App interface {
	Install(ctx context.Context, name, namespace string, version string, opt map[string]string) error
	Uninstall(ctx context.Context, name, namespace string) error
	Search(ctx context.Context, name string) (string, error)
}

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
		for _, a := range c.App.PreInstall.Apps {
			if err := Run(ctx, m, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
				return err
			}
		}
		// Run install
		if err := app.Install(ctx, appName, namespace, c.App.Version, nil); err != nil {
			return err
		}
		// Run postinstall
		for _, a := range c.App.PostInstall.Apps {
			if err := Run(ctx, m, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
				return err
			}
		}
	case Uninstall:
		return app.Uninstall(ctx, appName, namespace)
	default:
		return errors.New(fmt.Sprintf("Unsupported method %s", m))
	}
	return nil
}

func Search(args []string, configFile, namespace string) error {
	ctx := context.Background()
	c, err := config.New(configFile)
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

	if len(args) == 0 {
		out, err := app.Search(ctx, "")
		fmt.Print(string(out))
		return err
	}
	out, err := app.Search(ctx, args[0])
	fmt.Print(string(out))
	return err
}
