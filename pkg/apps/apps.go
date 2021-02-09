package apps

import (
	"context"
)

type App interface {
	Manifest(ctx context.Context, opt map[string]string) ([]byte, error)
}
