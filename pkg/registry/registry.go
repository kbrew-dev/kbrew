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

package registry

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/kbrew-dev/kbrew/pkg/config"
)

const (
	registriesDirName = "registries"

	defaultRegistryUserName = "kbrew-dev"
	defaultRegistryRepoName = "kbrew-registry"
	ghRegistryURLFormat     = "https://github.com/%s/%s.git"
	kbrewDir                = ".kbrew"
)

// recipeFilenamePattern regex pattern to search recipe files within a registry
var recipeFilenamePattern = regexp.MustCompile(`(?m)recipes\/(.*)\.yaml`)

// KbrewRegistry is the collection of kbrew recipes
type KbrewRegistry struct {
	path string
}

// Info holds recipe name and path for an app
type Info struct {
	Name string
	Path string
}

// New initializes KbrewRegistry, creates or clones default registry if not exists
func New(configDir string) (*KbrewRegistry, error) {
	r := &KbrewRegistry{
		path: filepath.Join(configDir),
	}
	return r, r.init()
}

// init creates config dir and clones default registry if not exists
func (kr *KbrewRegistry) init() error {
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	registriesDir := filepath.Join(home, kbrewDir, registriesDirName)

	// Check if default kbrew-registry exists, clone if not added already
	if _, err := os.Stat(filepath.Join(registriesDir, defaultRegistryUserName, defaultRegistryRepoName)); os.IsNotExist(err) {
		return kr.Add(defaultRegistryUserName, defaultRegistryRepoName, registriesDir)
	}
	return nil
}

// Add clones given kbrew registry in the config dir
func (kr *KbrewRegistry) Add(user, repo, path string) error {
	fmt.Printf("Adding %s/%s registry to %s\n", user, repo, path)
	r, err := git.PlainClone(filepath.Join(path, user, repo), false, &git.CloneOptions{
		URL:               fmt.Sprintf(ghRegistryURLFormat, user, repo),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return err
	}
	head, err := r.Head()
	if err != nil {
		return err
	}
	fmt.Printf("Registry %s/%s head at %s\n", user, repo, head)
	return err
}

// FetchRecipe iterates over all the kbrew recipes and returns path of the app recipe file
func (kr *KbrewRegistry) FetchRecipe(appName string) (string, error) {
	// Iterate over all the registries
	info, err := kr.Search(appName, true)
	if err != nil {
		return "", err
	}
	if len(info) == 0 {
		return "", fmt.Errorf("no recipe found for %s", appName)
	}
	return info[0].Path, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (kr *KbrewRegistry) FetchDetailRecipe(appName string) (string, error) {
	// Iterate over all the registries
	info, err := kr.Search(appName, true)
	if err != nil {
		return "", err
	}
	if len(info) == 0 {
		return "", fmt.Errorf("no recipe found for %s", appName)
	}
	 path := info[0].Path
	// fs,e os.OpenFile(path)
	data, err := os.ReadFile(path)
    check(err)
    //fmt.Print(string(dat))
	return string(data), nil
}

// Search returns app Info for give app
func (kr *KbrewRegistry) Search(appName string, exactMatch bool) ([]Info, error) {
	result := []Info{}
	appList, err := kr.ListApps()
	if err != nil {
		return nil, err
	}
	for _, app := range appList {
		if exactMatch {
			if app.Name == appName {
				return []Info{app}, nil
			}
			continue
		}
		if strings.HasPrefix(app.Name, appName) {
			result = append(result, app)
		}
	}
	return result, nil
}

// ListApps return Info list of all the apps
func (kr *KbrewRegistry) ListApps() ([]Info, error) {
	infoList := []Info{}
	err := filepath.WalkDir(kr.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir

		}
		if d.IsDir() {
			return nil
		}
		for _, match := range recipeFilenamePattern.FindAllStringSubmatch(path, -1) {
			if len(match) != 2 {
				continue
			}
			infoList = append(infoList, Info{Name: match[1], Path: path})
		}
		return nil
	})
	return infoList, err
}

// List returns list of registries
func (kr *KbrewRegistry) List() ([]string, error) {
	registries := []string{}

	// Registries are placed at - CONFIG_DIR/GITHUB_USER/GITHUB_REPO path
	// Interate over all the GITHUB_USERS dirs to find the list of all kbrew registries
	dirs, err := ioutil.ReadDir(kr.path)
	if err != nil {
		return nil, err
	}
	for _, user := range dirs {
		if !user.IsDir() {
			continue
		}
		subDirs, err := ioutil.ReadDir(filepath.Join(kr.path, user.Name()))
		if err != nil {
			return nil, err
		}
		for _, repo := range subDirs {
			if !repo.IsDir() {
				continue
			}
			registries = append(registries, fmt.Sprintf("%s/%s", user.Name(), repo.Name()))
		}
	}
	return registries, nil
}

// Update pull latest commits from registry repos
func (kr *KbrewRegistry) Update() error {
	registries, err := kr.List()
	if err != nil {
		return err
	}
	for _, r := range registries {
		if err := fetchUpdates(kr.path, r); err != nil {
			return err
		}
	}
	return nil
}

// Info returns information about a recipe
func (kr *KbrewRegistry) Info(appName string) (string, error) {
	c, err := kr.FetchRecipe(appName)
	if err != nil {
		return "", err
	}
	a, err := config.NewApp(appName, c)
	if err != nil {
		return "", err
	}
	bytes, err := yaml.Marshal(buildAppInfo(a.App))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Args returns the arguments declared for a recipe
func (kr *KbrewRegistry) Args(appName string) (map[string]interface{}, error) {
	c, err := kr.FetchRecipe(appName)
	if err != nil {
		return nil, err
	}
	a, err := config.NewApp(appName, c)
	if err != nil {
		return nil, err
	}

	return a.App.Args, nil
}

func fetchUpdates(rootDir, repo string) error {
	gitRegistry, err := git.PlainOpen(filepath.Join(rootDir, repo))
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return nil
		}
		return errors.Wrapf(err, "failed to init git repo %s", repo)
	}
	fmt.Printf("Fetching updates for registry %s\n", repo)
	wt, err := gitRegistry.Worktree()
	if err != nil {
		return errors.Wrapf(err, "failed to fetch updates for %s repo", repo)
	}
	err = wt.Pull(&git.PullOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Wrapf(err, "failed to fetch updates for %s repo", repo)
	}
	head, err := gitRegistry.Head()
	if err != nil {
		return errors.Wrapf(err, "failed to find head of %s repo", repo)
	}
	fmt.Printf("Registry %s head is set to %s\n", repo, head)
	return nil
}

func buildAppInfo(a config.App) config.App {
	app := config.App{
		Version: a.Version,
		Args:    a.Args,
		Repository: config.Repository{
			Name: a.Repository.Name,
			Type: a.Repository.Type,
			URL:  a.Repository.URL,
		},
	}
	preinstalls := []config.PreInstall{}
	for _, p := range a.PreInstall {
		if len(p.Apps) != 0 {
			preinstalls = append(preinstalls, config.PreInstall{
				Apps: p.Apps,
			})
		}
	}
	app.PreInstall = preinstalls
	postinstalls := []config.PostInstall{}
	for _, p := range a.PostInstall {
		if len(p.Apps) != 0 {
			postinstalls = append(postinstalls, config.PostInstall{
				Apps: p.Apps,
			})
		}
	}
	app.PostInstall = postinstalls
	return app
}
