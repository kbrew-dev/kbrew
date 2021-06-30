package apps

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"

	"github.com/kbrew-dev/kbrew/pkg/apps/helm"
	"github.com/kbrew-dev/kbrew/pkg/apps/raw"
	"github.com/kbrew-dev/kbrew/pkg/config"
	"github.com/kbrew-dev/kbrew/pkg/events"
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
	Workloads(ctx context.Context, namespace string) ([]corev1.ObjectReference, error)
}

// Run fetches recipe from registry for the app and performs given operation
func Run(ctx context.Context, m Method, appName, namespace, appConfigPath string) error {
	c, err := config.NewApp(appName, appConfigPath)
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
		return runInstall(ctx, app, c, m, appName, namespace, appConfigPath)
	case Uninstall:
		return runUninstall(ctx, app, c, m, appName, namespace, appConfigPath)
	default:
		return errors.New(fmt.Sprintf("Unsupported method %s", m))
	}
	return nil
}

func runInstall(ctx context.Context, app App, c *config.AppConfig, m Method, appName, namespace, appConfigPath string) error {
	// Event report
	event := events.NewKbrewEvent(c)

	// Run preinstall
	for _, phase := range c.App.PreInstall {
		for _, a := range phase.Apps {
			if err := Run(ctx, m, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
				return handleInstallError(ctx, err, event, app, appName, namespace)
			}
		}
		for _, a := range phase.Steps {
			if err := execCommand(ctx, a); err != nil {
				return handleInstallError(ctx, err, event, app, appName, namespace)
			}
		}
	}
	// Run install
	if err := app.Install(ctx, appName, namespace, c.App.Version, nil); err != nil {
		return handleInstallError(ctx, err, event, app, appName, namespace)
	}
	// Run postinstall
	for _, phase := range c.App.PostInstall {
		for _, a := range phase.Apps {
			if err := Run(ctx, m, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
				return handleInstallError(ctx, err, event, app, appName, namespace)
			}
		}
		for _, a := range phase.Steps {
			if err := execCommand(ctx, a); err != nil {
				return handleInstallError(ctx, err, event, app, appName, namespace)
			}
		}
	}
	if viper.GetBool(config.AnalyticsEnabled) {
		if err1 := event.Report(context.TODO(), events.ECInstallSuccess, nil, nil); err1 != nil {
			fmt.Printf("DEBUG: Failed to report event. %s\n", err1.Error())
		}
	}
	fmt.Printf("App %s installed successfully!\n", appName)
	return nil
}

func runUninstall(ctx context.Context, app App, c *config.AppConfig, m Method, appName, namespace, appConfigPath string) error {
	// Event report
	event := events.NewKbrewEvent(c)

	// Execute precleanup steps
	for _, a := range c.App.PreCleanup.Steps {
		if err := execCommand(ctx, a); err != nil {
			return handleUninstallError(ctx, err, event, appName, namespace)
		}
	}

	// Delete postinstall apps
	for _, phase := range c.App.PostInstall {
		for _, a := range phase.Apps {
			if err := Run(ctx, m, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
				return handleUninstallError(ctx, err, event, appName, namespace)
			}
		}
	}

	// Run uninstall
	if err := app.Uninstall(ctx, appName, namespace); err != nil {
		return handleUninstallError(ctx, err, event, appName, namespace)
	}

	// Delete preinstall apps
	for _, phase := range c.App.PreInstall {
		for _, a := range phase.Apps {
			if err := Run(ctx, m, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
				return handleUninstallError(ctx, err, event, appName, namespace)
			}
		}
	}

	// Execute postcleanup steps
	for _, a := range c.App.PostCleanup.Steps {
		if err := execCommand(ctx, a); err != nil {
			return handleUninstallError(ctx, err, event, appName, namespace)
		}
	}

	if viper.GetBool(config.AnalyticsEnabled) {
		if err1 := event.Report(context.TODO(), events.ECUninstallSuccess, nil, nil); err1 != nil {
			fmt.Printf("DEBUG: Failed to report event. %s\n", err1.Error())
		}
	}
	fmt.Printf("App %s uninstalled successfully!\n", appName)
	return nil
}

func handleInstallError(ctx context.Context, err error, event *events.KbrewEvent, app App, appName, namespace string) error {
	if err == nil {
		return nil
	}
	if !viper.GetBool(config.AnalyticsEnabled) {
		return err
	}
	wkl, err1 := app.Workloads(context.TODO(), namespace)
	if err1 != nil {
		fmt.Printf("DEBUG: Failed to report event. %s\n", err.Error())
	}

	if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
		fmt.Printf("ERROR: Timed out while installing %s app in %s namespace\n", appName, namespace)
		if err1 := event.Report(context.TODO(), events.ECInstallTimeout, err, nil); err1 != nil {
			fmt.Printf("DEBUG: Failed to report event. %s\n", err1.Error())
		}
		if err1 := event.ReportK8sEvents(context.TODO(), err, wkl); err1 != nil {
			fmt.Printf("DEBUG: Failed to report event. %s\n", err1.Error())
		}
		return err
	}
	if err1 := event.Report(context.TODO(), events.ECInstallFail, err, nil); err1 != nil {
		fmt.Printf("DEBUG: Failed to report event. %s\n", err1.Error())
	}
	if err1 := event.ReportK8sEvents(context.TODO(), err, wkl); err1 != nil {
		fmt.Printf("DEBUG: Failed to report event. %s\n", err1.Error())
	}
	return err
}

func handleUninstallError(ctx context.Context, err error, event *events.KbrewEvent, appName, namespace string) error {
	if err == nil {
		return nil
	}
	fmt.Printf("Warning: error encountered while uninstalling app - %s.\nYou need to cleanup few resources manually. App: %s, Namespace: %s\n", err, appName, namespace)
	if !viper.GetBool(config.AnalyticsEnabled) {
		return err
	}

	if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
		if err1 := event.Report(context.TODO(), events.ECUninstallTimeout, err, nil); err1 != nil {
			fmt.Printf("DEBUG: Failed to report event. %s\n", err1.Error())
		}
		return err
	}
	if err1 := event.Report(context.TODO(), events.ECUninstallFail, err, nil); err1 != nil {
		fmt.Printf("DEBUG: Failed to report event. %s\n", err1.Error())
	}
	return err
}

func execCommand(ctx context.Context, cmd string) error {
	c := exec.CommandContext(ctx, "sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
