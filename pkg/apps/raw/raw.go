package raw

import (
	"context"
	"os"
	"os/exec"

	"github.com/vishal-biyani/kbrew/pkg/apps"
	"github.com/vishal-biyani/kbrew/pkg/config"
)

type method string

const (
	install   method = "create"
	uninstall method = "delete"
	upgrade   method = "apply"
)

type RawApp struct {
	apps.BaseApp
}

func New(c config.App, namespace string) *RawApp {
	return &RawApp{
		apps.BaseApp{
			Name:      c.Name,
			Namespace: namespace,
			Version:   c.Version,
			URL:       c.URL,
			Digest:    c.SHA256,
		},
	}
}

func (r *RawApp) Install(ctx context.Context, options map[string]string) error {
	// TODO(@prasad): Use go sdks
	return kubectlCommand(install, r.Name, r.Namespace, r.URL)
}

func (r *RawApp) Uninstall(ctx context.Context) error {
	// TODO(@prasad): Use go sdks
	return kubectlCommand(uninstall, r.Name, r.Namespace, r.URL)
}

func kubectlCommand(m method, name, namespace, url string) error {
	c := exec.Command("kubectl", string(m), "-f", url, "--namespace", namespace)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
