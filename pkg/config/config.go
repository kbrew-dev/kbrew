package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	Raw  AppType = "raw"
	Helm AppType = "helm"
)

type AppType string

type AppConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	App        App    `yaml:"app"`
}

type App struct {
	Name    string  `yaml:"name"`
	Type    AppType `yaml:"type"`
	URL     string  `yaml:"url"`
	SHA256  string  `yaml:"sha256"`
	Version string  `yaml:"version"`
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
