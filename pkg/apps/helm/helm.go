package helm

import (
	"context"
	"os"
	"os/exec"

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
			Name:      c.Name,
			Namespace: namespace,
			Version:   c.Version,
			URL:       c.URL,
			Digest:    c.SHA256,
		},
	}
}

func (ha *HelmApp) Install(ctx context.Context, options map[string]string) error {
	//TODO: Resolve Deps
	// Validate and install chart
	// TODO(@prasad): Use go sdks
	return helmCommand(install, ha.Name, ha.Namespace, ha.URL)
}

func (ha *HelmApp) Uninstall(ctx context.Context) error {
	//TODO: Resolve Deps
	// Validate and install chart
	// TODO(@prasad): Use go sdks
	return helmCommand(uninstall, ha.Name, ha.Namespace, "")
}

func (ha *HelmApp) Manifest(ctx context.Context, opt map[string]string) ([]byte, error) {
	return nil, nil
}

func helmCommand(m method, name, namespace, url string) error {
	// Needs helm3
	c := exec.Command("helm", string(m), name, "--namespace", namespace, url)
	if url == "" {
		c = exec.Command("helm", string(m), name, "--namespace", namespace)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
