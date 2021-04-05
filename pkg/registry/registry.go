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
	"github.com/pkg/errors"
)

const (
	registriesDirName = "registries"

	defaultRegistryUserName = "kbrew-dev"
	defaultRegistryRepoName = "kbrew-registry"
	ghRegistryURLFormat     = "https://github.com/%s/%s.git"
)

// recipeFilenamePattern regex pattern to search recipe files within a registry
var recipeFilenamePattern = regexp.MustCompile(`(?m)recipes\/(.*)\.yaml`)

// kbrewRegistry is the collection of kbrew recipes
type kbrewRegistry struct {
	path string
}

// Info holds recipe name and path for an app
type Info struct {
	Name string
	Path string
}

// New initializes kbrewRegistry, creates or clones default registry if not exists
func New(configDir string) (*kbrewRegistry, error) {
	r := &kbrewRegistry{
		path: filepath.Join(configDir, registriesDirName),
	}
	return r, r.init()
}

// init creates config dir and clones default registry if not exists
func (kr *kbrewRegistry) init() error {
	// Generate registry path
	if _, err := os.Stat(kr.path); os.IsNotExist(err) {
		if err := os.MkdirAll(kr.path, os.ModePerm); err != nil {
			return errors.Wrap(err, "failed to initialize kbrew registry")
		}
	}

	// Check if default kbrew-registry exists, clone if not added already
	if _, err := os.Stat(filepath.Join(kr.path, defaultRegistryUserName, defaultRegistryRepoName)); os.IsNotExist(err) {
		return kr.Add(defaultRegistryUserName, defaultRegistryRepoName)
	}
	return nil
}

// Add clones given kbrew registry in the config dir
func (kr *kbrewRegistry) Add(user, repo string) error {
	fmt.Printf("Adding %s/%s registry to %s\n", user, repo, kr.path)
	r, err := git.PlainClone(filepath.Join(kr.path, user, repo), false, &git.CloneOptions{
		URL:               fmt.Sprintf(ghRegistryURLFormat, user, repo),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	head, err := r.Head()
	if err != nil {
		return err
	}
	fmt.Printf("Registry %s/%s head at %s\n", user, repo, head)
	return err
}

// FetchRecipe iterates over all the kbrew recipes and returns path of the app recipe file
func (kr *kbrewRegistry) FetchRecipe(appName string) (string, error) {
	// Iterate over all the registries
	info, err := kr.Search(appName, true)
	if err != nil {
		return "", err
	}
	if len(info) == 0 {
		return "", errors.New(fmt.Sprintf("No recipe found for %s", appName))
	}
	return info[0].Path, nil
}

// Search returns app Info for give app
func (kr *kbrewRegistry) Search(appName string, exactMatch bool) ([]Info, error) {
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
func (kr *kbrewRegistry) ListApps() ([]Info, error) {
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
			infoList = append(infoList, Info{Name: string(match[1]), Path: path})
		}
		return nil
	})
	return infoList, err
}

// List returns list of registries
func (kr *kbrewRegistry) List() ([]string, error) {
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
func (kr *kbrewRegistry) Update() error {
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
