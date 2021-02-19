package helm

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/vishal-biyani/kbrew/pkg/apps"
	"github.com/vishal-biyani/kbrew/pkg/config"
)

type method string

const (
	install   method = "install"
	uninstall method = "delete"
	upgrade   method = "upgrade"
)

type HelmApp struct {
	apps.BaseApp
}

func New(c config.App, namespace string) *HelmApp {
	return &HelmApp{
		apps.BaseApp{
			App:       c,
			Namespace: namespace,
		},
	}
}

func (ha *HelmApp) Install(ctx context.Context, name, version string, options map[string]string) error {
	fmt.Printf("Installing helm app %s/%s\n", ha.App.Repository.Name, name)
	//TODO: Resolve Deps
	// Validate and install chart
	// TODO(@prasad): Use go sdks
	// Needs helm3
	if _, err := ha.addRepo(ctx); err != nil {
		return err
	}
	out, err := helmCommand(install, name, version, ha.Namespace, fmt.Sprintf("%s/%s", ha.App.Repository.Name, name))
	fmt.Println(out)
	return err
}

func (ha *HelmApp) Uninstall(ctx context.Context, name string) error {
	fmt.Printf("Unistalling helm app %s\n", name)
	//TODO: Resolve Deps
	// Validate and install chart
	// TODO(@prasad): Use go sdks
	out, err := helmCommand(uninstall, name, "", ha.Namespace, "")
	fmt.Println(out)
	return err
}

func (ha *HelmApp) addRepo(ctx context.Context) (string, error) {
	// Needs helm3
	c := exec.Command("helm", "repo", "add", ha.App.Repository.Name, ha.App.Repository.URL)
	if out, err := c.CombinedOutput(); err != nil {
		return string(out), err
	}
	return ha.updateRepo(ctx)
}

func (ha *HelmApp) updateRepo(ctx context.Context) (string, error) {
	// Needs helm3
	c := exec.Command("helm", "repo", "update")
	out, err := c.CombinedOutput()
	return string(out), err
}

func (ha *HelmApp) Search(ctx context.Context, name string) (string, error) {
	// Needs helm3
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

func helmCommand(m method, name, version, namespace, chart string) (string, error) {
	// Needs helm3
	c := exec.Command("helm", string(m), name, "--namespace", namespace)
	if chart != "" {
		c.Args = append(c.Args, chart)
	}
	if version != "" {
		c.Args = append(c.Args, "--version", version)
	}
	out, err := c.CombinedOutput()
	return string(out), err
}
