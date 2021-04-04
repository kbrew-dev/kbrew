package update

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kbrew-dev/kbrew/pkg/version"

	"github.com/google/go-github/v27/github"
	"github.com/pkg/errors"
)

const (
	releaseRepoOwner = "kbrew-dev"
	releaseRepoName  = "kbrew-release"
	upgradeCmd       = "curl -sfL https://raw.githubusercontent.com/kbrew-dev/kbrew-release/main/install.sh | sh"
)

func getBinDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(path), nil
}

func CheckRelease(ctx context.Context) error {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(ctx, releaseRepoOwner, releaseRepoName)
	if err != nil {
		return errors.Wrap(err, "failed to check for kbrew updates")
	}
	if release == nil || release.TagName == nil {
		return nil
	}
	// Send notification if newer version available
	if version.VERSION != *release.TagName {
		fmt.Printf("kbrew %s is available, upgrading...\n", *release.TagName)
		return upgradeKbrew(ctx)
	}
	return nil
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
