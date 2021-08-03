// Copyright 2021 The kbrew Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

const (
	releaseRepoOwner = "kbrew-dev"
	releaseRepoName  = "kbrew"
)

// GetLatestVersion returns latest published release version on GitHub
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
