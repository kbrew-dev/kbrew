package util

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

const (
	releaseRepoOwner = "kbrew-dev"
	releaseRepoName  = "kbrew-release"
)

func GetLatestVersion(ctx context.Context) (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(ctx, releaseRepoOwner, releaseRepoName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check for kbrew updates")
	}
	if release == nil || release.TagName == nil {
		return nil, errors.Errorf("")
	}
	return release, nil
}
