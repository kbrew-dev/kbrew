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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

var (
	sampleYaml = `apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-scraping
`
	badYaml = `apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
	name: allow-scraping
`
	apiVersionYaml = `apiVersion: networking.k8s.io/v1beta1
kind: NetworkPolicy
metadata:
  name: allow-scraping
`
	rootExpression = "."

	apiVersionExpression = `select(.kind  == "NetworkPolicy" and .metadata.name == "allow-scraping").apiVersion |= "networking.k8s.io/v1beta1"`
)

func TestEval(t *testing.T) {

	type arg struct {
		manifest   string
		expression string
	}

	type want struct {
		result string
		err    error
	}

	cases := map[string]struct {
		arg
		want
	}{
		"CheckSimpleYaml": {
			arg: arg{
				manifest:   sampleYaml,
				expression: rootExpression,
			},
			want: want{
				result: sampleYaml,
			},
		},
		"CheckApiVersionChange": {
			arg: arg{
				manifest:   sampleYaml,
				expression: apiVersionExpression,
			},
			want: want{
				result: apiVersionYaml,
			},
		},
		"CheckBadYaml": {
			arg: arg{
				manifest:   badYaml,
				expression: rootExpression,
			},
			want: want{
				err: errors.Wrap(errors.New("yaml: line 4: found character that cannot start any token"), "Failed to evaluate"),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := NewEvaluator()
			o, err := e.Eval(tc.arg.manifest, tc.arg.expression)

			if diff := cmp.Diff(tc.want.result, o); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}

			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected '%s', got %q", tc.err.Error(), err.Error())
			}
		})
	}
}
