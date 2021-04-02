package update

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/kbrew-dev/kbrew/pkg/version"

	"github.com/google/go-github/v27/github"
)

const (
	releaseRepoOwner = "kbrew-dev"
	releaseRepoName  = "kbrew-release"
	upgradeCmd       = "curl -sfL https://raw.githubusercontent.com/kbrew-dev/kbrew-release/main/install.sh | sh"
)

func CheckRelease(ctx context.Context) error {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(ctx, releaseRepoOwner, releaseRepoName)
	if err != nil {
		fmt.Println("Failed to check for kbrew updates.")
		return err
	}
	if release == nil || release.TagName == nil {
		return nil
	}
	// Send notification if newer version available
	if version.VERSION != *release.TagName {
		fmt.Printf("kbrew %s is available, upgrading...\n", version.VERSION)
		return upgradeKbrew(ctx)
	}
	return nil
}

func upgradeKbrew(ctx context.Context) error {
	return execCommand(ctx, upgradeCmd)
}

func execCommand(ctx context.Context, cmd string) error {
	c := exec.CommandContext(ctx, "sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
