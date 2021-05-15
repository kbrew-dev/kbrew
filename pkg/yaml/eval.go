package yaml

import (
	"bytes"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/pkg/errors"
	logging "gopkg.in/op/go-logging.v1"
)

type evaluator struct {
	s yqlib.StreamEvaluator
}

func NewEvaluator() evaluator {
	// TODO(@sahil-lakhwani): check if there's a better way of avoiding yq debug logs
	logging.SetLevel(logging.CRITICAL, "yq-lib")
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
