package raw

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"text/tabwriter"

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
			App:       c,
			Namespace: namespace,
		},
	}
}

func (r *RawApp) Install(ctx context.Context, name, version string, options map[string]string) error {
	// TODO(@prasad): Use go sdks
	return kubectlCommand(install, name, r.Namespace, r.App.Repository.URL)
}

func (r *RawApp) Uninstall(ctx context.Context, name string) error {
	// TODO(@prasad): Use go sdks
	return kubectlCommand(uninstall, name, r.Namespace, r.App.Repository.URL)
}

func (ha *RawApp) Search(ctx context.Context, name string) (string, error) {
	return printList(ha.App), nil
}

func kubectlCommand(m method, name, namespace, url string) error {
	c := exec.Command("kubectl", string(m), "-f", url, "--namespace", namespace)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func printList(app config.App) string {
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "NAME\tVERSION\tTYPE")
	fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", app.Name, app.Version, app.Repository.Type))
	w.Flush()
	return b.String()
}
