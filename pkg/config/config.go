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

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kbrew-dev/kbrew/pkg/engine"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/tools/clientcmd"
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
	Args        map[string]interface{} `yaml:"args,omitempty"`
	Repository  Repository             `yaml:"repository"`
	Name        string                 `yaml:"name,omitempty"`
	Namespace   string                 `yaml:"namespace,omitempty"`
	URL         string                 `yaml:"url,omitempty"`
	SHA256      string                 `yaml:"sha256,omitempty"`
	Version     string                 `yaml:"version,omitempty"`
	PreInstall  []PreInstall           `yaml:"pre_install,omitempty"`
	PostInstall []PostInstall          `yaml:"post_install,omitempty"`
	PreCleanup  AppCleanup             `yaml:"pre_cleanup,omitempty"`
	PostCleanup AppCleanup             `yaml:"post_cleanup,omitempty"`
}

// Repository is the repo for kbrew app
type Repository struct {
	Name string   `yaml:"name"`
	URL  string   `yaml:"url"`
	Type RepoType `yaml:"type"`
}

// PreInstall contains Apps and Steps that need to be installed/executed before installing the main app
type PreInstall struct {
	Apps  []string `yaml:"apps,omitempty"`
	Steps []string `yaml:"steps,omitempty"`
}

// PostInstall contains Apps and Steps that need to be installed/executed after installing the main app
type PostInstall struct {
	Apps  []string `yaml:"apps,omitempty"`
	Steps []string `yaml:"steps,omitempty"`
}

// AppCleanup contains steps to be executed before uninstalling applications
type AppCleanup struct {
	Steps []string `yaml:"steps,omitempty"`
}

// NewApp parses kbrew recipe configuration and returns AppConfig instance
func NewApp(name, path string) (*AppConfig, error) {
	c := &AppConfig{}
	configFile, err := os.Open(path)
	defer func() {
		err := configFile.Close()
		if err != nil {
			fmt.Printf("Error closing file :%v", err)
		}
	}()

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	k8sconfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to load Kubernetes config")
	}

	e := engine.NewEngine(k8sconfig)
	v, err := e.Render(string(b))
	if err != nil {
		return nil, err
	}

	if len(b) != 0 {
		err = yaml.Unmarshal([]byte(v), c)
		if err != nil {
			return nil, err
		}
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
