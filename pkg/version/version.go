package version

import (
	"context"
	"fmt"

	"github.com/kbrew-dev/kbrew/pkg/util"
	"go.etcd.io/etcd/version"
)

// Version The below variables are overridden using the build process
// name of the release
var Version = "dev"

// GitCommitID git commit id of the release
var GitCommitID = "none"

// BuildDate date for the release
var BuildDate = "unknown"

const versionLongFmt = `{"Version": "%s", "GitCommit": "%s", "BuildDate": "%s"}`

// Long long version of the release
func Long(ctx context.Context) string {
	release, err := util.GetLatestVersion(ctx)
	if err != nil {
		fmt.Printf("Error getting latest version of kbrew from Github: %s", err)
	} else {
		if version.Version != *release.TagName {
			fmt.Printf("There is a new version of kbrew available: %s, please run 'kbrew update' command to upgrade\n", *release.TagName)

		}
	}
	return fmt.Sprintf(versionLongFmt, Version, GitCommitID, BuildDate)
}

// Short short version of the release
func Short() string {
	return Version
}
