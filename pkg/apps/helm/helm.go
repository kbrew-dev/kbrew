package helm

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	"github.com/kbrew-dev/kbrew/pkg/apps/raw"
	"github.com/kbrew-dev/kbrew/pkg/config"
	"github.com/kbrew-dev/kbrew/pkg/log"
)

type method string

const (
	installMethod   method = "install"
	statusMethod    method = "status"
	uninstallMethod method = "delete"
	// getManifestMethod method = "get manifest" // unused
)

// App holds helm app details
type App struct {
	app config.App
	log *log.Logger
}

// New returns Helm App
func New(c config.App, log *log.Logger) *App {
	return &App{
		app: c,
		log: log,
	}
}

// Install installs the application specified by name, version and namespace.
func (ha *App) Install(ctx context.Context, name, namespace, version string, options map[string]string) error {
	//TODO: Resolve Deps
	// Validate and install chart
	// TODO(@prasad): Use go sdks
	// Needs helm3
	if _, err := ha.addRepo(ctx); err != nil {
		return err
	}
	if err := ha.resolveArgs(); err != nil {
		return err
	}
	_, err := helmCommand(ctx, statusMethod, name, "", namespace, "", nil)
	if err == nil {
		// helm release already exists, return from here
		ha.log.Warnf("helm app %s/%s already exists in %s namespace. Skipping...\n", ha.app.Repository.Name, name, namespace)
		return nil
	}

	out, err := helmCommand(ctx, installMethod, name, version, namespace, fmt.Sprintf("%s/%s", ha.app.Repository.Name, name), ha.app.Args)
	ha.log.Debug(out)
	return err
}

// Uninstall uninstalls the application specified by name and namespace.
func (ha *App) Uninstall(ctx context.Context, name, namespace string) error {
	//TODO: Resolve Deps
	// Validate and install chart
	// TODO(@prasad): Use go sdks
	out, err := helmCommand(ctx, uninstallMethod, name, "", namespace, "", nil)
	ha.log.Debug(out)
	return err
}

func (ha *App) resolveArgs() error {
	if len(ha.app.Args) != 0 {
		for arg, value := range ha.app.Args {
			if value == nil {
				ha.app.Args[arg] = ""
			}
		}
	}
	return nil
}

func (ha *App) addRepo(ctx context.Context) (string, error) {
	// Needs helm 3.2+
	c := exec.CommandContext(ctx, "helm", "repo", "add", ha.app.Repository.Name, ha.app.Repository.URL)
	if out, err := c.CombinedOutput(); err != nil {
		return string(out), err
	}
	return ha.updateRepo(ctx)
}

func (ha *App) updateRepo(ctx context.Context) (string, error) {
	// Needs helm 3.2+
	c := exec.CommandContext(ctx, "helm", "repo", "update")
	out, err := c.CombinedOutput()
	return string(out), err
}

func (ha *App) getManifests(ctx context.Context, namespace string) (string, error) {
	c := exec.CommandContext(ctx, "helm", "get", "manifest", ha.app.Name, "--namespace", namespace)
	out, err := c.CombinedOutput()
	return string(out), err
}

// Search searches the name passed in helm repo
func (ha *App) Search(ctx context.Context, name string) (string, error) {
	// Needs helm 3.2+
	if out, err := ha.addRepo(ctx); err != nil {
		return out, err
	}
	c := exec.CommandContext(ctx, "helm", "search", "repo", fmt.Sprintf("%s/%s", ha.app.Repository.Name, name))
	out, err := c.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	if strings.Contains(string(out), "No results found") {
		return string(out), errors.New("No results found")
	}
	return string(out), err
}

// Workloads returns K8s workload object reference list for the helm app
func (ha *App) Workloads(ctx context.Context, namespace string) ([]corev1.ObjectReference, error) {
	manifest, err := ha.getManifests(ctx, namespace)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get helm chart manifests")
	}
	return raw.ParseManifestYAML(manifest, namespace)
}

func helmCommand(ctx context.Context, m method, name, version, namespace, chart string, chartArgs map[string]interface{}) (string, error) {
	// Needs helm 3.2+
	c := exec.CommandContext(ctx, "helm", string(m), name, "--namespace", namespace)
	if chart != "" {
		c.Args = append(c.Args, chart)
	}
	if version != "" {
		c.Args = append(c.Args, "--version", version)
	}
	if m == installMethod {
		// Add extra time to wait arg so that context will be timeout out before helm command failure
		// This is for catching timeout through context instead of parsing helm command output
		// This might change once we switch to SDKs
		c.Args = append(c.Args, "--wait", "--timeout", "5h0m", "--create-namespace")
	}

	if len(chartArgs) != 0 {
		c.Args = append(c.Args, appendChartArgs(chartArgs)...)
	}

	out, err := c.CombinedOutput()
	return string(out), err
}

func appendChartArgs(args map[string]interface{}) []string {
	var s []string
	for k, v := range args {
		s = append(s, "--set", k+"="+fmt.Sprintf("%s", v))
	}
	return s
}
