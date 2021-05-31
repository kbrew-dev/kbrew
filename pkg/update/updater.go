package update

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kbrew-dev/kbrew/pkg/util"
	"github.com/kbrew-dev/kbrew/pkg/version"

	"github.com/pkg/errors"
)

const (
	upgradeCmd = "curl -sfL https://raw.githubusercontent.com/kbrew-dev/kbrew-release/main/install.sh | sh"
)

func getBinDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(path), nil
}

// IsAvailable checks if a new version of GitHub release available
func IsAvailable(ctx context.Context) (string, error) {
	release, err := util.GetLatestVersion(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to check for kbrew updates")
	}
	if version.Version != *release.TagName {
		return *release.TagName, nil
	}
	return "", nil
}

// CheckRelease checks for the latest release
func CheckRelease(ctx context.Context) error {
	release, err := IsAvailable(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to check for kbrew updates")
	}
	if release == "" {
		return nil
	}
	fmt.Printf("kbrew %s is available, upgrading...\n", release)
	return upgradeKbrew(ctx)
}

func upgradeKbrew(ctx context.Context) error {
	dir, err := getBinDir()
	if err != nil {
		return errors.Wrap(err, "failed to get executable dir")
	}
	os.Setenv("BINDIR", dir)
	defer os.Unsetenv("BINDIR")
	return execCommand(ctx, upgradeCmd)
}

func execCommand(ctx context.Context, cmd string) error {
	c := exec.CommandContext(ctx, "sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
