package raw

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TODO: relative path
const dataDir = "./data/raw"

type Raw struct {
	app     string
	options map[string]string
}

func New(name string) (*Raw, error) {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", dataDir, name)); os.IsNotExist(err) {
		return nil, errors.New("Unsupported application")
	}
	return &Raw{
		app: name,
	}, nil
}

func (r *Raw) Manifest(ctx context.Context, opt map[string]string) ([]byte, error) {
	return ParseRaw(fmt.Sprintf("./data/raw/%s", r.app))
}

func List(ctx context.Context) ([]string, error) {
	appList := []string{}
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return appList, err
	}

	for _, file := range files {
		appList = append(appList, file.Name())
	}
	return appList, nil
}

// TODO: Move to utils
func ParseRaw(path string) ([]byte, error) {
	resp := []byte{}
	yamlFiles, err := filepath.Glob(path + "/*.yaml")
	if err != nil {
		return nil, err
	}
	for _, f := range yamlFiles {
		file, err := os.Open(f) // For read access.
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		resp = append(resp, []byte("---\n")...)
		resp = append(resp, b...)
	}
	return resp, nil
}
