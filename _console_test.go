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
	"testing"
)

func TestConsole(t *testing.T) {
	config := `{"console":{"Severity":5, "DebugLevel":3}}`
	l := NewLogDeb(10, config)
    defer l.Destroy()
	l.Deb("TestConsole", "test console")
	l.Debl("TestConsole", "test console with debug level 4", 4) // Don't write this because DebugLevel is set to 3
	l.Debl("TestConsole", "test console with debug level 3", 3)

	config = `{"console":{"Severity":5}}`
	l2 := NewLogDeb(10, config)
    defer l2.Destroy()
	l2.Debl("TestConsole2", "test console with debug level 1", 1)
	l2.Debl("TestConsole2", "test console with debug level 2", 2)

    config = `{"main":{"UseFncRules":true},"console":{"Severity":5, "FncRules":{"TestConsole3.writeme":{"Severity":5},"TestConsole3.dontwriteme":{"Severity":2}}}}`
	l3 := NewLogDeb(10, config)
    defer l3.Destroy()
	l3.Deb("TestConsole3.writeme", "test console with rule 1")
	l3.Deb("TestConsole3.dontwriteme", "test console with rule 2")

	config = `{"console":{}}`
	l4 := NewLogDeb(10, config)
    defer l4.Destroy()
	l4.Deb("TestConsole4.dontwriteme", "test console debug without config")

}

func TestConsoleMainConf(t *testing.T) {
	config := `{"main":{"Severity":5}, "console":{}}`
	l5 := NewLogDeb(10, config)
    defer l5.Destroy()
	l5.Deb("TestConsole5.writeme", "Debug with MAIN severity Debug")

	config = `{"main":{"Severity":5}, "console": {}}`
	l6 := NewLogDeb(10, config)
    defer l6.Destroy()
	l6.Err("TestConsole6.writeme", "Error with MAIN severity Debug")

	config = `{"main":{"Severity":2}, "console": {}}`
	l7 := NewLogDeb(10, config)
    defer l7.Destroy()
	l7.Warn("TestConsole7.dontwriteme", "Warning with MAIN severity Error")

    config = `{"main":{"Severity":2}, "console":{"FncRules":{"TestConsole8.dontwriteme":{"Severity":5}}}}`
	l8 := NewLogDeb(10, config)
    defer l8.Destroy()
    l8.Warn("TestConsole8.dontwriteme", "Warning with MAIN severity Error. FncRule but don't use it")

    config = `{"main":{"Severity":2, "UseFncRules":true}, "console":{"FncRules":{"TestConsole9.writeme":{"Severity":5}}}}`
	l9 := NewLogDeb(10, config)
    defer l9.Destroy()
    l9.Warn("TestConsole9.writeme", "Warning with MAIN severity Error. But use FncRule")
}
