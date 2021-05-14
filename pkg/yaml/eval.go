package yaml

import (
	"bytes"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/pkg/errors"
)

type evaluator struct {
	s yqlib.StreamEvaluator
}

func NewEvaluator() evaluator {
	return evaluator{s: yqlib.NewStreamEvaluator()}
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
