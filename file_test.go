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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func runTestFile(t *testing.T, name string, config string, tmsgs []STLogMsg) {
	var unconf map[string]interface{}
	var filename string
	err := json.Unmarshal([]byte(config), &unconf)
	if err != nil {
		panic(fmt.Sprintf("Error extracting writer config from json. ERR: %s", err))
	}
	if fc, t := unconf["file"]; t {
		fcm := fc.(map[string]interface{})
		if fn, t := fcm["filename"]; t {
			filename = fn.(string)
		} else {
			panic(fmt.Sprintf("File name non configured for File Writer"))
		}
	} else {
		panic(fmt.Sprintf("Can not find configuration for File Writer"))
	}
	// Delete the log file to get a clean config
	os.Remove(filename)
	executeTest(config, tmsgs)
	f, err := os.OpenFile(filename, os.O_RDONLY, 0660)
	if err != nil {
		fmt.Println(err)
	}
	outb, err := ioutil.ReadAll(bufio.NewReader(f))
	out := string(outb)
	prTest("FILE OUTPUT:", out)
	checkResult(t, out, name, FILESEP, tmsgs)
}

func TestFileDeb1(t *testing.T) {
	name := "TestFileDeb1"
	fnc := tFncName(name)
	config := `{"file":{"flags":0, "filename":"file.log"}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: fnc, msg: "test file", sev: SEVERROR}, true},
	}
	runTestFile(t, name, config, tmsgs[:])
}

func TestFileDeb2(t *testing.T) {
	name := "TestFileDeb2"
	fnc := tFncName(name)
	config := `{"file":{"flags":0, "sev":5, "dlev":3, "filename":"file.log"}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: fnc, msg: "test file"}, true},
	}
	runTestFile(t, name, config, tmsgs[:])
}

func TestFileDeb3(t *testing.T) {
	name := "TestFileDeb3"
	fnc := tFncName(name)
	config := `{"file":{"flags":0, "sev":5, "dlev":3, "filename":"file.log"}}`
	tmsgs := []STLogMsg{
		STLogMsg{SLogMsg{fnc: fnc, msg: "test file 1"}, true},
		STLogMsg{SLogMsg{fnc: fnc, msg: "test file 2"}, true},
	}
	runTestFile(t, name, config, tmsgs[:])
}
