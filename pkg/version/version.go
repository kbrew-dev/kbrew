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

package version

import "fmt"

// Version The below variables are overridden using the build process
// name of the release
var Version = "dev"

// GitCommitID git commit id of the release
var GitCommitID = "none"

// BuildDate date for the release
var BuildDate = "unknown"

const versionLongFmt = `{"Version": "%s", "GitCommit": "%s", "BuildDate": "%s"}`

// Long long version of the release
func Long() string {
	return fmt.Sprintf(versionLongFmt, Version, GitCommitID, BuildDate)
}

// Short short version of the release
func Short() string {
	return Version
}
