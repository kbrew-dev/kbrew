package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	// Analytics setting flags

	// AnalyticsUUID represents unique ID used as customer ID
	AnalyticsUUID = "analyticsUUID"
	// AnalyticsEnabled to toggle GA event collection
	AnalyticsEnabled = "analyticsEnabled"
)

// KbrewConfig is a kbrew config stored at CONFIG_DIR/config.yaml
type KbrewConfig struct {
	AnalyticsUUID    string `yaml:"analyticsUUID"`
	AnalyticsEnabled bool   `yaml:"analyticsEnabled"`
}

// AppConfig is the kbrew recipe configuration
type AppConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	App        App    `yaml:"app"`
}

// App hold app details set in kbrew recipe
type App struct {
	Args        map[string]interface{} `yaml:"args"`
	Repository  Repository             `yaml:"repository"`
	Name        string                 `yaml:"name"`
	Namespace   string                 `yaml:"namespace"`
	URL         string                 `yaml:"url"`
	SHA256      string                 `yaml:"sha256"`
	Version     string                 `yaml:"version"`
	PreInstall  []PreInstall           `yaml:"pre_install"`
	PostInstall []PostInstall          `yaml:"post_install"`
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

// NewApp parses kbrew recipe configuration and returns AppConfig instance
func NewApp(name, path string) (*AppConfig, error) {
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
	c.App.Name = name
	return c, nil
}

// NewKbrew parses Kbrew config and returns KbrewConfig struct object
func NewKbrew() (*KbrewConfig, error) {
	kc := &KbrewConfig{}
	err := viper.Unmarshal(kc)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read kbrew config")
	}
	return kc, nil
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

	// Create kbrew config yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigDir)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create config file
			viper.Set(AnalyticsUUID, uuid.NewV4().String())
			viper.Set(AnalyticsEnabled, true)
			cobra.CheckErr(viper.SafeWriteConfig())
		} else {
			cobra.CheckErr(err)
		}
	}
}
