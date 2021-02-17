package apps

import (
	"context"
)

type App interface {
	Install(ctx context.Context, opt map[string]string) error
	Uninstall(ctx context.Context) error
}

type BaseApp struct {
	URL       string
	Name      string
	Namespace string
	Digest    string
	Version   string
}
