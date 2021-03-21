package registry

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
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
		// TODO(@prasad): Remove auth once registry is made public
		Auth: &githttp.BasicAuth{
			Username: "PrasadG193",
			Password: os.Getenv("GITHUB_TOKEN"),
		},
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
	appList, err := kr.List()
	if err != nil {
		return nil, err
	}
	for _, app := range appList {
		if exactMatch && app.Name == appName {
			return []Info{app}, nil
		}
		if strings.HasPrefix(app.Name, appName) {
			result = append(result, app)
		}
	}
	return result, nil
}

// List return Info list of all the apps
func (kr *kbrewRegistry) List() ([]Info, error) {
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
