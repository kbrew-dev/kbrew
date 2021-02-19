package apps

import (
	"context"

	"github.com/vishal-biyani/kbrew/pkg/config"
)

type App interface {
	Install(ctx context.Context, name string, version string, opt map[string]string) error
	Uninstall(ctx context.Context, name string) error
	Search(ctx context.Context, name string) (string, error)
}

type BaseApp struct {
	App       config.App
	Namespace string
}
