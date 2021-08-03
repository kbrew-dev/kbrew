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

package yaml

import (
	"bytes"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/pkg/errors"
	logging "gopkg.in/op/go-logging.v1"
)

// Evaluator parses yaml documents and patches them by applying given expressions
type Evaluator interface {
	Eval(manifest string, expresssion string) (string, error)
}

type evaluator struct {
	s yqlib.StreamEvaluator
}

// NewEvaluator returns a new wrapped streamevaluator
func NewEvaluator() Evaluator {
	// TODO(@sahil-lakhwani): check if there's a better way of avoiding yq debug logs
	logging.SetLevel(logging.CRITICAL, "yq-lib")
	return &evaluator{s: yqlib.NewStreamEvaluator()}
}

func (e *evaluator) Eval(manifest string, expresssion string) (string, error) {

	node, err := yqlib.NewExpressionParser().ParseExpression(expresssion)
	if err != nil {
		return "", errors.Wrap(err, "Error evaluating expression")
	}

	reader := strings.NewReader(manifest)

	var buf bytes.Buffer
	printer := yqlib.NewPrinter(&buf, false, true, false, 2, true)

	err = e.s.Evaluate("", reader, node, printer)
	if err != nil {
		return "", errors.Wrap(err, "Failed to evaluate")
	}

	return buf.String(), nil
}
