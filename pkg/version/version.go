package version

import "fmt"

// The below variables are overrriden using the build process
var VERSION = "dev"
var GIT_COMMIT_ID = "none"
var BUILD_DATE = "unknown"

const versionLongFmt = `{"Version": "%s", "GitCommit": "%s", "BuildDate": "%s"}`

func Long() string {
	return fmt.Sprintf(versionLongFmt, VERSION, GIT_COMMIT_ID, BUILD_DATE)
}

func Short() string {
	return VERSION
}
