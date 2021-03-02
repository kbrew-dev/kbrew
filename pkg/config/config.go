package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	Raw  RepoType = "raw"
	Helm RepoType = "helm"
)

type RepoType string

type AppConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	App        App    `yaml:"app"`
}

type App struct {
	Repository  Repository    `yaml:"repository"`
	Name        string        `yaml:"name"`
	Namespace   string        `yaml:"namespace"`
	URL         string        `yaml:"url"`
	SHA256      string        `yaml:"sha256"`
	Version     string        `yaml:"version"`
	PreInstall  []PreInstall  `yaml:"pre_install"`
	PostInstall []PostInstall `yaml:"post_install"`
}

type Repository struct {
	Name string   `yaml:"name"`
	URL  string   `yaml:"url"`
	Type RepoType `yaml:"type"`
}

type PreInstall struct {
	Apps  []string
	Steps []string
}

type PostInstall struct {
	Apps  []string
	Steps []string
}

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
