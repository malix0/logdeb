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
	"bytes"
	"io"
	"os"
	"testing"
)

type STLogMsg struct {
	SLogMsg
	logit bool
}

var oldStout *os.File
var w *os.File
var outC chan string

func preTestConsole() {
	var r *os.File
	oldStout = os.Stdout // keep backup of the real stdout
	r, w, _ = os.Pipe()
	os.Stdout = w
	outC = make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
}

func postTestConsole() {
	// back to normal state
	w.Close()
	os.Stdout = oldStout // restoring the real stdout
}

// Run console writer test and capture output from stdout
func runTestConsole(t *testing.T, name string, config string, tmsgs []STLogMsg) {
	preTestConsole()
	executeTest(config, tmsgs)
	postTestConsole()
	out := <-outC
	prTest("CONSOLE OUTPUT:", out)
	checkResult(t, out, name, CONSSEP, tmsgs)
}

func TestConsoleDeb1(t *testing.T) {
	name := "TestConsoleDeb1"
	fnc := tFncName(name)
	config := `{"console":{}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: fnc, msg: "test console debug without config"}, false},
	}
	runTestConsole(t, name, config, tmsgs[:])
}

func TestConsoleDeb2(t *testing.T) {
	name := "TestConsoleDeb2"
	fnc := tFncName(name)
	config := `{"console":{"flags":0, "sev":5, "dlev":3}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: fnc, msg: "test console"}, true},
		// Don't write this because the configured DebugLevel is set to 3
		STLogMsg{SLogMsg{fnc: fnc, msg: "test console with debug level 4", debLev: 4}, false},
		STLogMsg{SLogMsg{fnc: fnc, msg: "test console with debug level 3", debLev: 3}, true},
	}
	runTestConsole(t, name, config, tmsgs[:])
}

func TestConsoleDeb3(t *testing.T) {
	name := "TestConsoleDeb3"
	fnc := tFncName(name)
	config := `{"console":{"flags":0, "sev":5}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: fnc, msg: "test console with debug level 1", debLev: 1}, true},
		// Don't write this because the default DebugLevel is 1
		STLogMsg{SLogMsg{fnc: fnc, msg: "test console with debug level 2", debLev: 2}, false},
	}
	runTestConsole(t, name, config, tmsgs[:])
}

func TestConsoleDeb4(t *testing.T) {
	name := "TestConsoleDeb4"
	config := `{"main":{"usefncrules":true},"console":{"flags":0, "sev":5, "fncrules":{"TestConsoleDeb4.writeme":{"sev":5},"TestConsoleDeb4.dontwriteme":{"sev":2}}}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: "TestConsoleDeb4.writeme", msg: "test console with rule 1"}, true},
		// Don't write this because the function severity is set to 2
		STLogMsg{SLogMsg{fnc: "TestConsoleDeb4.dontwriteme", msg: "test console with rule 2"}, false},
	}
	runTestConsole(t, name, config, tmsgs[:])
}

func TestConsoleMainConf1(t *testing.T) {
	name := "TestConsoleMainConf1"
	fnc := tFncName(name)
	config := `{"main":{"sev":5}, "console":{"flags":0}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: fnc, msg: "Debug with MAIN severity Debug"}, true},
		STLogMsg{SLogMsg{fnc: fnc, sev: SEVERROR, msg: "Error with MAIN severity Debug"}, true},
	}
	runTestConsole(t, name, config, tmsgs[:])
}

func TestConsoleMainConf2(t *testing.T) {
	name := "TestConsoleMainConf2"
	fnc := tFncName(name)
	config := `{"main":{"sev":2}, "console":{"flags":0}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: fnc, sev: SEVWARN, msg: "Warning with MAIN severity Error"}, false},
	}
	runTestConsole(t, name, config, tmsgs[:])
}

func TestConsoleMainConf3(t *testing.T) {
	name := "TestConsoleMainConf3"
	config := `{"main":{"sev":2}, "console":{"flags":0, "fncrules":{"TestConsoleMainConf3.dontwriteme":{"sev":5}}}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: "TestConsoleMainConf3.dontwriteme", sev: SEVWARN, msg: "Warning with MAIN severity Error. FncRule but don't use it"}, false},
	}
	runTestConsole(t, name, config, tmsgs[:])
}

func TestConsoleMainConf4(t *testing.T) {
	name := "TestConsoleMainConf4"
	config := `{"main":{"sev":2, "usefncrules":true}, "console":{"flags":0, "fncrules":{"TestConsoleMainConf4.writeme":{"sev":3}}}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: "TestConsoleMainConf4.writeme", sev: SEVWARN, msg: "Warning with MAIN severity Error. But use FncRule"}, true},
	}
	runTestConsole(t, name, config, tmsgs[:])
}
