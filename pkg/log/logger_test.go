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

package log

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type LoggerTestSuite struct{}

var _ = Suite(&LoggerTestSuite{})

func (l *LoggerTestSuite) TestLogger(c *C) {
	var buf bytes.Buffer
	log := NewLogger(true)
	log.SetWriter(&buf)

	log.Info("info log", "app=app1")
	c.Assert(buf.String(), Equals, "\rinfo log app=app1\n")

	buf.Reset()
	log.Infof("info log app=%s, namespace=%s", "app1", "appns")
	c.Assert(buf.String(), Equals, "\rinfo log app=app1, namespace=appns\n")

	buf.Reset()
	log.Debug("debug log", "app=app1")
	c.Assert(buf.String(), Equals, "\rDEBUG: debug log app=app1\n")

	buf.Reset()
	log.Debugf("debug log app=%s, namespace=%s", "app1", "appns")
	c.Assert(buf.String(), Equals, "\rDEBUG: debug log app=app1, namespace=appns\n")

	buf.Reset()
	log.Warn("warn log", "app=app1")
	c.Assert(buf.String(), Equals, "\r\x1b[33mWARN:\x1b[0m warn log app=app1\n")

	buf.Reset()
	log.Warnf("warn log app=%s, namespace=%s", "app1", "appns")
	c.Assert(buf.String(), Equals, "\r\x1b[33mWARN:\x1b[0m warn log app=app1, namespace=appns\n")

	buf.Reset()
	log.Error("error log", "app=app1")
	c.Assert(buf.String(), Equals, "\r\x1b[31mERROR:\x1b[0m error log app=app1\n")

	buf.Reset()
	log.Errorf("error log app=%s, namespace=%s", "app1", "appns")
	c.Assert(buf.String(), Equals, "\r\x1b[31mERROR:\x1b[0m error log app=app1, namespace=appns\n")

	buf.Reset()
	log.InfoMap("app", "app1")
	c.Assert(buf.String(), Equals, "\r\x1b[32mapp:\x1b[0m app1\n")

	buf.Reset()
	log.InfoMap("namespace", "appns")
	c.Assert(buf.String(), Equals, "\r\x1b[32mnamespace:\x1b[0m appns\n")
}
