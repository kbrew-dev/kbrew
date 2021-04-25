package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// RepoType describes the type of kbrew app repository
type RepoType string

// ConfigDir represents dir path of kbrew config
var ConfigDir string

const (
	// Raw repo type means the apps in the repo are raw apps
	Raw RepoType = "raw"
	// Helm repo means the apps in the repo are helm apps
	Helm RepoType = "helm"
	// RegistriesDirName represents the dir name within ConfigDir holding all the kbrew registries
	RegistriesDirName = "registries"
)

// AppConfig is the kbrew recipe configuration
type AppConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	App        App    `yaml:"app"`
}

// App hold app details set in kbrew recipe
type App struct {
	Args        map[string]string `yaml:"args"`
	Repository  Repository        `yaml:"repository"`
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	URL         string            `yaml:"url"`
	SHA256      string            `yaml:"sha256"`
	Version     string            `yaml:"version"`
	PreInstall  []PreInstall      `yaml:"pre_install"`
	PostInstall []PostInstall     `yaml:"post_install"`
}

// Repository is the repo for kbrew app
type Repository struct {
	Name string   `yaml:"name"`
	URL  string   `yaml:"url"`
	Type RepoType `yaml:"type"`
}

// PreInstall contains Apps and Steps that need to be installed/executed before installing the main app
type PreInstall struct {
	Apps  []string
	Steps []string
}

// PostInstall contains Apps and Steps that need to be installed/executed after installing the main app
type PostInstall struct {
	Apps  []string
	Steps []string
}

// New parses kbrew recipe configuration and returns AppConfig instance
func New(path string) (*AppConfig, error) {
	c := &AppConfig{}
	configFile, err := os.Open(path)
	defer configFile.Close()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	if len(b) != 0 {
		yaml.Unmarshal(b, c)
	}
	return c, nil
}

// InitConfig initializes ConfigDir.
// If ConfigDir does not exists, create it
func InitConfig() {
	if ConfigDir != "" {
		return
	}
	// Find home directory.
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	// Generate default config file path
	ConfigDir = filepath.Join(home, ".kbrew")
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		err := os.MkdirAll(ConfigDir, os.ModePerm)
		cobra.CheckErr(err)
	}
}
