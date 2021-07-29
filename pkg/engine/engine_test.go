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

package engine

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestFuncMap(t *testing.T) {
	fns := funcMap()
	forbidden := []string{"env", "expandenv"}
	for _, f := range forbidden {
		if _, ok := fns[f]; ok {
			t.Errorf("Forbidden function %s exists in FuncMap.", f)
		}
	}

	// Test for Engine-specific template functions.
	expect := []string{"include", "required", "tpl", "toYaml", "fromYaml", "toToml", "toJson", "fromJson", "lookup"}
	for _, f := range expect {
		if _, ok := fns[f]; !ok {
			t.Errorf("Expected add-on function %q", f)
		}
	}
}

func TestRender(t *testing.T) {
	type want struct {
		result string
		err    error
	}

	cases := map[string]struct {
		arg string
		want
	}{
		"CheckUntil": {
			arg: "{{ until 5 }}",
			want: want{
				result: "[0 1 2 3 4]",
			},
		},
		"CheckMissingFunction": {
			arg: "{{ foo 5 }}",
			want: want{
				err: errors.Wrapf(errors.New("template: gotpl:1: function \"foo\" not defined"), renderErr),
			},
		},
		"CheckConstString": {
			arg: "SomeString",
			want: want{
				result: "SomeString",
			},
		},
		"CheckInclude": {
			arg: `{{define "T1"}}{{trim .}}{{end}}{{include "T1" " hello" | upper }}`,
			want: want{
				result: "HELLO",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := NewEngine(nil)
			o, err := e.Render(tc.arg)

			if diff := cmp.Diff(tc.want.result, o); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}

			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected '%s', got %q", tc.err.Error(), err.Error())
			}
		})
	}
}

func TestInitMap(t *testing.T) {
	e := NewEngine(nil)
	e.initFuncMap()

	_, ok := e.fmap["include"]
	if !ok {
		t.Error("Expected function 'include' in funcMap")
	}
}
