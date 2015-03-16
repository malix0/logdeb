// Copyright 2014 Massimo Fidanza.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package logdeb

import (
	"fmt"
	"strings"
	"testing"
)

// Init logger and run tests
func executeTest(config string, tmsgs []STLogMsg) {
	l := NewLogDeb(10, config)
	for _, tm := range tmsgs {
		if tm.sev == 0 {
			tm.sev = SEVDEBUG
		}
		if tm.sev == SEVDEBUG {
			if tm.debLev != 0 {
				l.Debl(tm.fnc, tm.msg, tm.debLev)
			} else {
				l.Deb(tm.fnc, tm.msg)
			}
		} else if tm.sev == SEVINFO {
			l.Info(tm.fnc, tm.msg)
		} else if tm.sev == SEVWARN {
			l.Warn(tm.fnc, tm.msg)
		} else if tm.sev == SEVERROR {
			l.Err(tm.fnc, tm.msg)
		} else if tm.sev == SEVFATAL {
			l.Fatal(tm.fnc, tm.msg)
		}
	}
	// Destroy force the flush of the message channel
	l.Destroy()
}

// Check the test result
func checkResult(t *testing.T, out string, name string, sep string, tmsgs []STLogMsg) {
	outs := strings.Split(out, "\n")
	var found bool
	var matchlog, expect string
	err := false
	for _, tm := range tmsgs {
		found = false
		if tm.sev == 0 {
			tm.sev = SEVDEBUG
		}
		matchlog = fmt.Sprintf("%s[%s] %s %s", tm.fnc, tm.sev, sep, tm.msg)
		if expect != "" {
			expect = expect + "\n"
		}
		expect = expect + matchlog
		for i := 0; i < len(outs); i++ {
			//if strings.Contains(outs[i], matchlog) {
			if matchlog == outs[i] {
				found = true
				break
			}
		}
		if tm.logit != found {
			err = true
		}
	}
	if err {
		t.Errorf("%s\n EXPECT => %v\n GOT => %v", name, expect, out)
	}
}
