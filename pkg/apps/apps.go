package apps

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"

	"github.com/kbrew-dev/kbrew/pkg/apps/helm"
	"github.com/kbrew-dev/kbrew/pkg/apps/raw"
	"github.com/kbrew-dev/kbrew/pkg/config"
	"github.com/kbrew-dev/kbrew/pkg/events"
	"github.com/kbrew-dev/kbrew/pkg/log"
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

type AppRunner struct {
	operation Method
	log       *log.Logger
	status    *log.Status
}

func NewAppRunner(op Method, log *log.Logger, status *log.Status) *AppRunner {
	return &AppRunner{
		operation: op,
		log:       log,
		status:    status,
	}
}

// Run fetches recipe from registry for the app and performs given operation
func (r *AppRunner) Run(ctx context.Context, appName, namespace, appConfigPath string) error {
	c, err := config.NewApp(appName, appConfigPath)
	if err != nil {
		return err
	}
	var app App

	switch c.App.Repository.Type {
	case config.Helm:
		app = helm.New(c.App, r.log)
	case config.Raw:
		app, err = raw.New(c.App, r.log)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported app type %s", c.App.Repository.Type)
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

	switch r.operation {
	case Install:
		return r.runInstall(ctx, app, c, appName, namespace, appConfigPath)
	case Uninstall:
		return r.runUninstall(ctx, app, c, appName, namespace, appConfigPath)
	default:
		err = fmt.Errorf("unsupported method %s", r.operation)
	}
	return err
}

func (r *AppRunner) runInstall(ctx context.Context, app App, c *config.AppConfig, appName, namespace, appConfigPath string) error {
	// Event report
	event := events.NewKbrewEvent(c)

	// Run preinstall
	r.status.Start(fmt.Sprintf("Setting up pre-install dependencies for %s", appName))
	for _, phase := range c.App.PreInstall {
		for _, a := range phase.Apps {
			if err := r.Run(ctx, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
				return r.handleInstallError(ctx, err, event, app, appName, namespace)
			}
		}
		for _, a := range phase.Steps {
			out, err := r.execCommand(ctx, a)
			if err != nil {
				return r.handleInstallError(ctx, err, event, app, appName, namespace)
			}
			r.log.Debug(out)
		}
	}
	r.status.Stop()

	// Run install
	r.status.Start(fmt.Sprintf("Installing app %s in %s namespace", appName, namespace))
	if err := app.Install(ctx, appName, namespace, c.App.Version, nil); err != nil {
		return r.handleInstallError(ctx, err, event, app, appName, namespace)
	}
	r.status.Success()

	// Run postinstall
	r.status.Start(fmt.Sprintf("Setting up post-install dependencies for %s", appName))
	for _, phase := range c.App.PostInstall {
		for _, a := range phase.Apps {
			if err := r.Run(ctx, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
				return r.handleInstallError(ctx, err, event, app, appName, namespace)
			}
		}
		for _, a := range phase.Steps {
			out, err := r.execCommand(ctx, a)
			if err != nil {
				return r.handleInstallError(ctx, err, event, app, appName, namespace)
			}
			r.log.Debug(out)
		}
	}
	r.status.Stop()
	if viper.GetBool(config.AnalyticsEnabled) {
		if err1 := event.Report(context.TODO(), events.ECInstallSuccess, nil, nil); err1 != nil {
			r.log.Debugf("Failed to report event. %s\n", err1.Error())
		}
	}
	return nil
}

func (r *AppRunner) runUninstall(ctx context.Context, app App, c *config.AppConfig, appName, namespace, appConfigPath string) error {
	// Event report
	event := events.NewKbrewEvent(c)

	r.status.Start(fmt.Sprintf("Executing up pre-cleanup steps for %s", appName))
	// Execute precleanup steps
	for _, a := range c.App.PreCleanup.Steps {
		out, err := r.execCommand(ctx, a)
		if err != nil {
			return r.handleUninstallError(ctx, err, event, appName, namespace)
		}
		r.log.Debug(out)
	}
	r.status.Stop()

	// Delete postinstall apps
	for _, phase := range c.App.PostInstall {
		for _, a := range phase.Apps {
			if err := r.Run(ctx, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
				return r.handleUninstallError(ctx, err, event, appName, namespace)
			}
		}
	}

	// Run uninstall
	r.status.Start(fmt.Sprintf("Removing app %s from %s namespace", appName, namespace))
	if err := app.Uninstall(ctx, appName, namespace); err != nil {
		return r.handleUninstallError(ctx, err, event, appName, namespace)
	}
	r.status.Success()

	// Delete preinstall apps
	for _, phase := range c.App.PreInstall {
		for _, a := range phase.Apps {
			if err := r.Run(ctx, a, namespace, filepath.Join(filepath.Dir(appConfigPath), a+".yaml")); err != nil {
				return r.handleUninstallError(ctx, err, event, appName, namespace)
			}
		}
	}

	// Execute postcleanup steps
	r.status.Start(fmt.Sprintf("Executing up post-cleanup steps for %s", appName))
	for _, a := range c.App.PostCleanup.Steps {
		out, err := r.execCommand(ctx, a)
		if err != nil {
			return r.handleUninstallError(ctx, err, event, appName, namespace)
		}
		r.log.Debug(out)
	}
	r.status.Stop()

	if viper.GetBool(config.AnalyticsEnabled) {
		if err1 := event.Report(context.TODO(), events.ECUninstallSuccess, nil, nil); err1 != nil {
			r.log.Debugf("Failed to report event. %s\n", err1.Error())
		}
	}
	return nil
}

func (r *AppRunner) handleInstallError(ctx context.Context, err error, event *events.KbrewEvent, app App, appName, namespace string) error {
	if err == nil {
		return nil
	}
	defer r.status.Error()

	eventType := events.ECInstallFail
	if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
		r.log.Errorf("Timed out while installing %s app in %s namespace\n", appName, namespace)
		eventType = events.ECInstallTimeout
	}

	if !viper.GetBool(config.AnalyticsEnabled) {
		return err
	}

	wkl, err1 := app.Workloads(context.TODO(), namespace)
	if err1 != nil {
		r.log.Debugf("Failed to report event. %s\n", err.Error())
	}
	if err1 := event.Report(context.TODO(), eventType, err, nil); err1 != nil {
		r.log.Debugf("Failed to report event. %s\n", err1.Error())
	}
	if err1 := event.ReportK8sEvents(context.TODO(), err, wkl); err1 != nil {
		r.log.Debugf("Failed to report event. %s\n", err1.Error())
	}
	return err
}

func (r *AppRunner) handleUninstallError(ctx context.Context, err error, event *events.KbrewEvent, appName, namespace string) error {
	if err == nil {
		return nil
	}
	defer r.status.Error()
	r.log.Warnf("Error encountered while uninstalling app - %s.\nYou need to cleanup few resources manually. App: %s, Namespace: %s\n", err, appName, namespace)
	if !viper.GetBool(config.AnalyticsEnabled) {
		return err
	}

	if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
		if err1 := event.Report(context.TODO(), events.ECUninstallTimeout, err, nil); err1 != nil {
			r.log.Debugf("Failed to report event. %s\n", err1.Error())
		}
		return err
	}
	if err1 := event.Report(context.TODO(), events.ECUninstallFail, err, nil); err1 != nil {
		r.log.Debugf("Failed to report event. %s\n", err1.Error())
	}
	return err
}

func (r *AppRunner) execCommand(ctx context.Context, cmd string) (string, error) {
	c := exec.CommandContext(ctx, "sh", "-c", cmd)
	c.Stderr = os.Stderr
	output, err := c.Output()
	return string(output), err
}
