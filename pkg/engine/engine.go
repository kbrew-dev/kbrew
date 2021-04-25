package engine

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/engine"
	"k8s.io/client-go/rest"
)

// Engine provides functionalities to evaluate and render templates.
type Engine struct {
	fmap template.FuncMap
}

// NewEngine returns an initialized Engine
func NewEngine(config *rest.Config) *Engine {
	return &Engine{
		fmap: funcMap(config),
	}
}

// Render resolves values of a string template.
func (e *Engine) Render(arg string) (string, error) {

	keys := make([]string, 0, len(e.fmap))
	for k := range e.fmap {
		keys = append(keys, k)
	}

	t := template.Must(template.New("helm").Funcs(e.fmap).Parse(arg))

	var tpl bytes.Buffer
	err := t.Execute(&tpl, "")
	if err != nil {
		return "", errors.Wrapf(err, "Error rendering value")
	}

	return tpl.String(), nil
}

func funcMap(config *rest.Config) template.FuncMap {

	fmap := sprig.TxtFuncMap()

	fmap["lookup"] = engine.NewLookupFunction(config)

	return fmap
}
