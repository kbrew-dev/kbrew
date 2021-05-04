package helm

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/kbrew-dev/kbrew/pkg/config"
)

type method string

const (
	installMethod   method = "install"
	uninstallMethod method = "delete"
	upgrade         method = "upgrade"
)

// App holds helm app details
type App struct {
	App config.App
}

// New returns Helm App
func New(c config.App) *App {
	return &App{
		App: c,
	}
}

// Install installs the application specified by name, version and namespace.
func (ha *App) Install(ctx context.Context, name, namespace, version string, options map[string]string) error {
	fmt.Printf("Installing helm app %s/%s\n", ha.App.Repository.Name, name)
	//TODO: Resolve Deps
	// Validate and install chart
	// TODO(@prasad): Use go sdks
	// Needs helm3
	if _, err := ha.addRepo(ctx); err != nil {
		return err
	}

	out, err := helmCommand(installMethod, name, version, namespace, fmt.Sprintf("%s/%s", ha.App.Repository.Name, name), ha.App.Args)
	fmt.Println(out)
	return err
}

// Uninstall uninstalls the application specified by name and namespace.
func (ha *App) Uninstall(ctx context.Context, name, namespace string) error {
	fmt.Printf("Unistalling helm app %s\n", name)
	//TODO: Resolve Deps
	// Validate and install chart
	// TODO(@prasad): Use go sdks
	out, err := helmCommand(uninstallMethod, name, "", namespace, "", nil)
	fmt.Println(out)
	return err
}

func (ha *App) addRepo(ctx context.Context) (string, error) {
	// Needs helm 3.2+
	c := exec.Command("helm", "repo", "add", ha.App.Repository.Name, ha.App.Repository.URL)
	if out, err := c.CombinedOutput(); err != nil {
		return string(out), err
	}
	return ha.updateRepo(ctx)
}

func (ha *App) updateRepo(ctx context.Context) (string, error) {
	// Needs helm 3.2+
	c := exec.Command("helm", "repo", "update")
	out, err := c.CombinedOutput()
	return string(out), err
}

// Search searches the name passed in helm repo
func (ha *App) Search(ctx context.Context, name string) (string, error) {
	// Needs helm 3.2+
	if out, err := ha.addRepo(ctx); err != nil {
		return string(out), err
	}
	c := exec.Command("helm", "search", "repo", fmt.Sprintf("%s/%s", ha.App.Repository.Name, name))
	out, err := c.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	if strings.Contains(string(out), "No results found") {
		return string(out), errors.New("No results found")
	}
	return string(out), err
}

func helmCommand(m method, name, version, namespace, chart string, chartArgs map[string]string) (string, error) {
	// Needs helm 3.2+
	c := exec.Command("helm", string(m), name, "--namespace", namespace)
	if chart != "" {
		c.Args = append(c.Args, chart)
	}
	if version != "" {
		c.Args = append(c.Args, "--version", version)
	}
	if m == installMethod {
		c.Args = append(c.Args, "--wait", "--create-namespace")
	}

	if chartArgs != nil && len(chartArgs) != 0 {
		c.Args = append(c.Args, appendChartArgs(chartArgs)...)
	}

	out, err := c.CombinedOutput()
	return string(out), err
}

func appendChartArgs(args map[string]string) []string {
	var s []string
	for k, v := range args {
		s = append(s, "--set", k+"="+v)
	}
	return s
}
