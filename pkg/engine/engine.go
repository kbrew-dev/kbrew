package engine

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/engine"
	"k8s.io/client-go/rest"
)

const renderErr = "Error rendering value"

// Engine provides functionalities to evaluate and render templates.
type Engine struct {
	config   *rest.Config
	fmap     template.FuncMap
	template *template.Template
}

// NewEngine returns an initialized Engine
func NewEngine(config *rest.Config) *Engine {
	return &Engine{
		template: template.New("gotpl"),
		config:   config,
	}
}

// Render resolves values of a string template.
func (e *Engine) Render(arg string) (string, error) {

	if len(e.fmap) == 0 {
		e.initFuncMap()
	}

	_, err := e.template.Parse(arg)
	if err != nil {
		return "", errors.Wrapf(err, renderErr)
	}

	var tpl bytes.Buffer
	err = e.template.Execute(&tpl, "")
	if err != nil {
		return "", errors.Wrapf(err, renderErr)
	}

	return tpl.String(), nil
}

func (e *Engine) initFuncMap() {

	includedNames := make(map[string]int)

	e.fmap = funcMap()

	e.fmap["include"] = func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > 1000 {
				return "", errors.Wrapf(fmt.Errorf("unable to execute template"), "rendering template has a nested reference name: %s", name)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := e.template.ExecuteTemplate(&buf, name, data)
		includedNames[name]--
		return buf.String(), err
	}

	if e.config != nil {
		e.fmap["lookup"] = engine.NewLookupFunction(e.config)
	}

	e.template.Funcs(e.fmap)
}
